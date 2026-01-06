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
				"go run /code/main < /code/input.txt",
			},
		)

	default:
		return "", "", errors.New("unsupported language")
	}
}
