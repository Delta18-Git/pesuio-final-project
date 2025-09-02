package compiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/go-archive"
)

func runDocker(tempFile *os.File, language string, input string) (string, string) {
	var cli *client.Client

	var containerConfig *container.Config = &container.Config{
		Tty:       false,
		Cmd:       []string{"/bin/sh", "-c"},
		OpenStdin: true,
	}

	var hostConfig *container.HostConfig = &container.HostConfig{
		Resources: container.Resources{
			Memory:    250 * 1024 * 1024,
			CPUQuota:  50000,  //half a CPU
			CPUPeriod: 100000, //per second
		},
	}

	var containerImage string
	var compileCommand string
	fileName := filepath.Base(tempFile.Name())

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err.Error()
	}
	defer cli.Close()

	bgContext := context.Background()

	switch language {
	case "py":
		containerImage = "python:slim"
		compileCommand = "python3 /" + fileName
	case "go":
		containerImage = "golang:alpine"
		compileCommand = "go run /" + fileName
	case "c":
		containerImage = "frolvlad/alpine-gcc:latest"
		compileCommand = "gcc -static -o /compiledcode /" + fileName + "; ./compiledcode"
	case "cpp":
		containerImage = "frolvlad/alpine-gxx:latest"
		compileCommand = "g++ -static -o /compiledcode /" + fileName + "; ./compiledcode"
	}
	imageList, err := cli.ImageList(bgContext, image.ListOptions{})
	if err != nil {
		return "", fmt.Sprintf("failed to list Docker images: %v", err)
	}
	imageExists := false

	for _, image := range imageList {
		if slices.Contains(image.RepoTags, containerImage) {
			imageExists = true
		}
		if imageExists {
			break
		}
	}
	if !imageExists {
		pullOutput, err := cli.ImagePull(bgContext, containerImage, image.PullOptions{})
		if err != nil {
			return "", fmt.Sprintf("failed to get image from registry: %v", err)
		}
		defer pullOutput.Close()
		io.Copy(io.Discard, pullOutput) // Wait for pull to end before continuing.
	}
	containerConfig.Image = containerImage
	containerConfig.Cmd = append(containerConfig.Cmd, compileCommand)
	createdContainer, err := cli.ContainerCreate(bgContext, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", err.Error()
	}
	defer cli.ContainerRemove(bgContext, createdContainer.ID, container.RemoveOptions{
		RemoveVolumes: true,
	})
	tarTempFile, tarErr := archive.Tar(tempFile.Name(), archive.Uncompressed)
	if tarErr != nil {
		return "", fmt.Sprintf("failed to write code to tar file: %v", err)
	}
	copyErr := cli.CopyToContainer(bgContext, createdContainer.ID, "/", tarTempFile, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if copyErr != nil {
		return "", fmt.Sprintf("failed to copy file to container: %v", err)
	}

	hijackedResponse, err := cli.ContainerAttach(bgContext, createdContainer.ID, container.AttachOptions{
		Stdin:  true,
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return "", fmt.Sprintf("failed to write input to container: %v", err)
	}

	err = cli.ContainerStart(bgContext, createdContainer.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Sprintf("failed to start container: %v", err)
	}

	statusCh, errCh := cli.ContainerWait(bgContext, createdContainer.ID, container.WaitConditionNotRunning)
	var stdout, stderr bytes.Buffer
	go func() {
		_, _ = stdcopy.StdCopy(&stdout, &stderr, hijackedResponse.Reader)
	}()
	_, err = hijackedResponse.Conn.Write([]byte(input + "\n"))
	if err != nil {
		return "", fmt.Sprintf("failed to write input to stdin: %v", err)
	}
	hijackedResponse.CloseWrite()
	select {
	case <-time.After(2 * time.Minute):
		stopErr := cli.ContainerStop(context.Background(), createdContainer.ID, container.StopOptions{})
		if stopErr != nil {
			return "", fmt.Sprintf("program execution timed out, failed to stop container: %v", stopErr)
		}
		return "", "program execution timed out"
	case err := <-errCh:
		return "", fmt.Sprintf("error waiting for container: %v", err)
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return stdout.String(), fmt.Sprintf("status code: %d\n %v", status.StatusCode, stderr.String())
		}
	}
	return stdout.String(), stderr.String()
}
