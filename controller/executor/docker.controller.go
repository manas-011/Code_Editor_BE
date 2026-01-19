package executor

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"fmt"

	"github.com/manas-011/code-editor-backend/model"
)

func runDockerCmd(
	ctx context.Context,
	image string,
	tmpDir string,
	cmd []string,
) (stdout, stderr string, exitCode int) {

	args := []string{
		"run", "--rm",
		"--cpus=0.5",
		"--memory=256m",
		"--network=none",
		"-v", tmpDir + ":/code",
		"-w", "/code",
		image,
	}
	args = append(args, cmd...)

	c := exec.CommandContext(ctx, "docker", args...)

	var outBuf, errBuf bytes.Buffer
	c.Stdout = &outBuf
	c.Stderr = &errBuf

	err := c.Run()

	exitCode = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func RunInDocker(
	ctx context.Context,
	extension string,
	image string,
	code string,
	input string,
	compileCmd []string,
	runCmd []string,
) (*model.JudgeResult, error) {

	tmpDir, err := os.MkdirTemp("", "exec-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	_ = os.WriteFile(filepath.Join(tmpDir, "main."+extension), []byte(code), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "input.txt"), []byte(input), 0644)

	// ---------- COMPILE ----------
	if compileCmd != nil {
		_, compileStderr, exitCode := runDockerCmd(ctx, image, tmpDir, compileCmd)
		if exitCode != 0 {
			return &model.JudgeResult{
				Status:   "CE",
				Stderr:  compileStderr,
				ExitCode: exitCode,
			}, nil
		}
	}

	// ---------- RUN ----------
	runCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stdout, stderr, exitCode := runDockerCmd(runCtx, image, tmpDir, runCmd)

	if runCtx.Err() == context.DeadlineExceeded {
		return &model.JudgeResult{Status: "TLE"}, nil
	}

	// Docker OOM kill
	if exitCode == 137 {
		return &model.JudgeResult{Status: "MLE"}, nil
	}

	if exitCode != 0 {
		return &model.JudgeResult{
			Status:   "RE",
			Stdout:   stdout,
			Stderr:   stderr,
			ExitCode: exitCode,
		}, nil
	}

	fmt.Println(stdout);
	fmt.Println(stderr);

	return &model.JudgeResult{
		Status:   "success",
		Stdout:   stdout,
		Stderr:   stderr, // warnings live here
		ExitCode: 0,
	}, nil
}
