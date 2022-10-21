package goroutine

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
)

const (
	ErrCodeGoroutineNoLogger = "GORTN-001"
)

var (
	ErrGoroutineNoLogger = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeGoroutineNoLogger, "either logger or logger func must be specified").C(ctx).Err()
	}
)
