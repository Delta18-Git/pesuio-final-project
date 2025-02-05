package compiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

func runDocker(tempFile *os.File, language string, input string) (outputString, errorString string) {
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
		panic(err.Error())
	}
	defer cli.Close()

	var emptyContext context.Context = context.Background()

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
		containerImage = "frolvlad/alpine-gxx"
		compileCommand = "g++ -static -o /compiledcode /" + fileName + "; ./compiledcode"
	}
	imageList, err := cli.ImageList(emptyContext, image.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list Docker images: %w", err).Error()
	}
	imageExists := false
	for _, image := range imageList {
		for _, tag := range image.RepoTags {
			if tag == containerImage {
				imageExists = true
				break
			}
		}
		if imageExists {
			break
		}
	}
	if !imageExists {
		pullOutput, err := cli.ImagePull(emptyContext, containerImage, image.PullOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get image from registry: %v", err).Error()
		}
		defer pullOutput.Close()
		io.Copy(io.Discard, pullOutput) //Wait for pull to end before continuing.
	}
	containerConfig.Image = containerImage
	containerConfig.Cmd = append(containerConfig.Cmd, compileCommand)
	createdContainer, err := cli.ContainerCreate(emptyContext, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		panic(err.Error())
	}
	defer cli.ContainerRemove(emptyContext, createdContainer.ID, container.RemoveOptions{
		RemoveVolumes: true,
	})
	tarTempFile, tarErr := archive.Tar(tempFile.Name(), archive.Uncompressed)
	if tarErr != nil {
		return "", fmt.Errorf("failed to write code to tar file: %v", err).Error()
	}
	copyErr := cli.CopyToContainer(emptyContext, createdContainer.ID, "/", tarTempFile, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if copyErr != nil {
		return "", fmt.Errorf("failed to copy file to container: %v", err).Error()
	}

	hijackedResponse, err := cli.ContainerAttach(emptyContext, createdContainer.ID, container.AttachOptions{
		Stdin:  true,
		Stream: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to write input to container: %v", err).Error()
	}

	if err := cli.ContainerStart(emptyContext, createdContainer.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	_, err = hijackedResponse.Conn.Write([]byte(input + "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to write input to stdin: %v", err).Error()
	}
	defer hijackedResponse.Close()
	output, err := cli.ContainerLogs(emptyContext, createdContainer.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return "", fmt.Errorf("failed to write get output from logs: %v", err).Error()
	}

	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, output)
	outputString = string(stdout.Bytes())
	errorString = string(stderr.Bytes())

	return outputString, errorString
}
