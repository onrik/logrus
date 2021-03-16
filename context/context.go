package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

type logCtxKey struct{}

// WithLogger puts logger into a context.
func WithLogger(parent context.Context, l logrus.FieldLogger) context.Context {
	return context.WithValue(parent, logCtxKey{}, l)
}

// GetLogger returns a logger from context.
// If any logger isn't found, GetLogger returns standard instance of logrus.
func GetLogger(ctx context.Context) logrus.FieldLogger {
	l, ok := ctx.Value(logCtxKey{}).(logrus.FieldLogger)
	if !ok {
		logrus.Debug("logger is not found in a context")
		return logrus.StandardLogger()
	}

	return l
}
