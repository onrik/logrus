package filename

import (
	"fmt"
	"runtime"

	"github.com/Sirupsen/logrus"
)

type Hook struct {
	levels []logrus.Level
}

func NewHook(levels ...logrus.Level) *Hook {
	return &Hook{
		levels: levels,
	}
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data["filename"] = getCaller()
	return nil
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(5)
	if !ok {
		return ""
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return fmt.Sprintf("%s:%d", file, line)

}
