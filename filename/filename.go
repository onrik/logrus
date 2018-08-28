package filename

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Hook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
	// ModulesToIgnore allows specifying multiple package to skip
	// when looking for the correct line up the stack.
	// Usefult e.g. with a wrapper around logrus.
	ModulesToIgnore map[string]struct{}
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip, hook.ModulesToIgnore))
	return nil
}

func NewHook(levels ...logrus.Level) *Hook {
	hook := Hook{
		Field:  "source",
		Skip:   5,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
		ModulesToIgnore: map[string]struct{}{"logrus": struct{}{}},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook
}

func findCaller(skip int, modulesToIgnore map[string]struct{}) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)

	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		module := strings.Split(file, "/")[0]
		if _, ok := modulesToIgnore[module]; ok {
			continue
		}
		break
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
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

	return pc, file, line
}
