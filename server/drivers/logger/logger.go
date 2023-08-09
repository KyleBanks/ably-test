package logger

import "context"

type Logger interface {
	Info(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
}
