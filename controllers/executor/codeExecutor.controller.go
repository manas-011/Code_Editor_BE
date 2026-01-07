package executor

import (
	"context"
	"errors"
)

func Execute(ctx context.Context, lang, code, input string) (string, string, error) {
	switch lang {

	case "go":
		return RunInDocker(
			ctx,
			"golang:1.22-alpine",
			code,
			input,
			[]string{
				"sh", "-c",
				"go run /code/main.go < /code/input.txt",
			},
		)

	case "cpp":
		return RunInDocker(
			ctx,
			"gcc:13",
			code,
			input,
			[]string{
				"sh", "-c",
				"g++ /code/main.cpp -O2 -o /code/a.out && /code/a.out < /code/input.txt",
			},
		)

	case "java":
		return RunInDocker(
			ctx,
			"openjdk:21-slim",
			code,
			input,
			[]string{
				"sh", "-c",
				"javac /code/Main.java && java -cp /code Main < /code/input.txt",
			},
		)

	case "python":
		return RunInDocker(
			ctx,
			"python:3.12-slim",
			code,
			input,
			[]string{
				"sh", "-c",
				"python3 /code/main.py < /code/input.txt",
			},
		)

	default:
		return "", "", errors.New("unsupported language")
	}
}
