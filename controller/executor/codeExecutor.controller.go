package executor

import (
	"context"
	"errors"

	"github.com/manas-011/code-editor-backend/model"
)

func Execute(ctx context.Context, lang, code, input string) (*model.JudgeResult, error) {
	switch lang {

	case "go":
		return RunInDocker(
			ctx,
			"go",
			"golang:1.22-alpine",
			code,
			input,
			[]string{"go", "build", "-o", "app", "main.go"},
			[]string{"sh", "-c", "./app < input.txt"},
		)

	case "cpp":
		return RunInDocker(
			ctx,
			"cpp",
			"gcc:13",
			code,
			input,
			[]string{"g++", "main.cpp", "-O2", "-std=gnu++20", "-o", "a.out"},
			[]string{"sh", "-c", "./a.out < input.txt"},
		)

	case "java":
		return RunInDocker(
			ctx,
			"java",
			"eclipse-temurin:21-jdk",
			code,
			input,
			[]string{"javac", "Main.java"},
			[]string{"sh", "-c", "java Main < input.txt"},
		)

	case "python":
		return RunInDocker(
			ctx,
			"py",
			"python:3.12-slim",
			code,
			input,
			nil, // no compile step
			[]string{"sh", "-c", "python3 main.py < input.txt"},
		)

	default:
		return nil, errors.New("unsupported language")
	}
}
