package stdout

import (
	"context"
	"fmt"
	"os"
)

type Logger struct{}

func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, args...))
}

func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(msg, args...))
}
