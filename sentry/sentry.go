package sentry

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

var (
	levelsMap = map[logrus.Level]raven.Severity{
		logrus.PanicLevel: raven.FATAL,
		logrus.FatalLevel: raven.FATAL,
		logrus.ErrorLevel: raven.ERROR,
		logrus.WarnLevel:  raven.WARNING,
		logrus.InfoLevel:  raven.INFO,
		logrus.DebugLevel: raven.DEBUG,
	}
)

const (
	HttpRequestField = "http_request"
	StacktaceField   = "stacktrace"
)

type Hook struct {
	client           *raven.Client
	levels           []logrus.Level
	HttpRequestField string
	StacktaceField   string
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	culprit := ""
	interfaces := []raven.Interface{&raven.Message{entry.Message, nil}}
	if err, ok := entry.Data[logrus.ErrorKey].(error); ok {
		culprit = err.Error()
		entry.Data[logrus.ErrorKey] = culprit
	}
	if httpRequest, ok := entry.Data[hook.HttpRequestField].(*http.Request); ok {
		interfaces = append(interfaces, raven.NewHttp(httpRequest))
		delete(entry.Data, hook.HttpRequestField)
	}

	if stacktrace, ok := entry.Data[hook.StacktaceField].(*raven.Stacktrace); ok {
		interfaces = append(interfaces, stacktrace)
		delete(entry.Data, hook.StacktaceField)
	}

	packet := &raven.Packet{
		Message:    entry.Message,
		Level:      levelsMap[entry.Level],
		Interfaces: interfaces,
		Extra:      entry.Data,
		Culprit:    culprit,
	}

	_, ch := hook.client.Capture(packet, map[string]string{})
	if entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel {
		return <-ch
	}

	return nil
}

func (hook *Hook) SetTags(tags map[string]string) {
	hook.client.Tags = tags
}

func (hook *Hook) AddTag(key, value string) {
	hook.client.Tags[key] = value
}

func (hook *Hook) SetRelease(release string) {
	hook.client.SetRelease(release)
}

func (hook *Hook) SetEnvironment(environment string) {
	hook.client.SetEnvironment(environment)
}

func NewHook(dsn string, levels ...logrus.Level) *Hook {
	client, err := raven.New(dsn)
	if err != nil {
		logrus.WithError(err).Error("Set DSN error")
	}

	hook := Hook{
		client:           client,
		levels:           levels,
		HttpRequestField: HttpRequestField,
		StacktaceField:   StacktaceField,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook
}
