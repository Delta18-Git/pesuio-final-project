package compiler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	var request models.RunRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	tempFile, err := os.CreateTemp("", "code-*."+request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file to store code"})
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(request.Code))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write code to temp file"})
		return
	}
	out, err := func() (string, error) {
		var cmd *exec.Cmd
		switch request.Language {
		case "py":
			{
				cmd = exec.Command("python3", (*os.File)(tempFile).Name())
			}
		case "go":
			{
				cmd = exec.Command("go", "run", (*os.File)(tempFile).Name())
			}
		case "c":
			{
				cmd = exec.Command("gcc", "-o", "unsafe-code", (*os.File)(tempFile).Name())
				_, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to compile C code: %v", err)
				}
				cmd = exec.Command("./unsafe-code")
				defer os.Remove("./unsafe-code")
			}
		case "cpp":
			{
				cmd = exec.Command("g++", "-o", "unsafe-code", (*os.File)(tempFile).Name())
				_, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to compile C++ code: %v", err)
				}
				cmd = exec.Command("./unsafe-code")
				defer os.Remove("./unsafe-code")
			}
		default:
			{
				return "", fmt.Errorf("unsupported language")
			}
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}
		defer stdin.Close()
		_, err = io.WriteString(stdin, string(request.Input)+"\n")
		if err != nil {
			return "", fmt.Errorf("failed to write input to stdin: %v", err)
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		} else {
			return string(out), nil
		}
	}()
	var output string
	if err != nil {
		output = err.Error()
	} else {
		output = out
	}
	c.JSON(http.StatusOK, gin.H{
		"output": output,
	})
}
