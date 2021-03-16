package context

import (
	"bytes"
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestWithLogger(t *testing.T) {
	l := logrus.New()
	b := bytes.NewBufferString("")
	l.SetOutput(b)

	GetLogger(WithLogger(context.Background(), l)).Info("hi")

	require.NotEmpty(t, b.String())
}

func TestGetLogger(t *testing.T) {
	require.NotNil(t, GetLogger(context.Background()))
}
