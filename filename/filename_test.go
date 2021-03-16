package filename

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestFilenameHook(t *testing.T) {
	hook := NewHook()

	buff := new(bytes.Buffer)
	logrus.SetOutput(buff)
	logrus.AddHook(hook)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Info("Test")

	require.NotEqual(t, "", buff.String())
	require.Equal(t,
		"filename/filename_test.go:20",
		gjson.Get(buff.String(), "_source").Str,
	)
}
