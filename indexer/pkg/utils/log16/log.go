package log16

import (
	"context"
	"github.com/inconshreveable/log15"
)

type Logger struct {
	log15.Logger
}

func NewLogger(key, val string, opts ...Option) *Logger {
	return &Logger{log15.New(key, val)}
}

func (l *Logger) Critical(ctx context.Context, msg string, ctxs ...interface{}) {
	i := []interface{}{"stack", string(Stack(2))}
	i = append(i, ctxs...)
	l.Crit(msg, i...)
}

func (l *Logger) Err(ctx context.Context, msg string, ctxs ...interface{}) {
	//ctxs = append(ctxs, "stack", Stack(2))
	i := []interface{}{"stack", string(Stack(2))}
	i = append(i, ctxs...)
	l.Error(msg, i...)
}
