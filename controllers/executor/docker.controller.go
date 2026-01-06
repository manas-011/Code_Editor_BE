package executor

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

func RunInDocker(
	ctx context.Context,
	image string,
	code string,
	input string,
	runCmd []string,
) (stdout string, stderr string, err error) {

	tmpDir, err := os.MkdirTemp("", "exec-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	codePath := filepath.Join(tmpDir, "main")
	inputPath := filepath.Join(tmpDir, "input.txt")

	if err = os.WriteFile(codePath, []byte(code), 0644); err != nil {
		return "", "", err
	}
	if err = os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		return "", "", err
	}

	args := []string{
		"run", "--rm",
		"--cpus=0.5",
		"--memory=256m",
		"--network=none",
		"-v", tmpDir + ":/code",
		image,
	}
	args = append(args, runCmd...)

	cmd := exec.CommandContext(ctx, "docker", args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}
