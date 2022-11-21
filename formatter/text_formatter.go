package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37

	defaultSourceField = "_source"
)

const defaultTimestampFormat = time.RFC3339

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	logrus.TextFormatter
	isTerminal       bool
	terminalInitOnce sync.Once

	SourceField  string
	FormatValue  func(value interface{}) string
	FormatCaller func(caller *runtime.Frame) string
}

// New creates new formatter
func New() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
	if f.FormatValue == nil {
		f.FormatValue = DefaultFormat
	}
	if f.FormatCaller == nil {
		f.FormatCaller = CallerFormat
	}
	if f.SourceField == "" {
		f.SourceField = defaultSourceField
	}
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.terminalInitOnce.Do(func() { f.init(entry) })

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		}
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message)
		}
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}
	if entry.HasCaller() && entry.Data["_source"] == nil {
		f.appendKeyValue(b, f.SourceField, f.FormatCaller(entry.Caller))
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	b.WriteString(f.FormatValue(value))
}

// DefaultFormat format value as text with fmt.Sprintf
func DefaultFormat(value interface{}) string {
	s, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%+v", value)
	}

	if len(s) == 0 || NeedsQuoting(s) {
		return fmt.Sprintf("%q", s)
	}

	return s
}

// JSONFormat format value as json
func JSONFormat(value interface{}) string {
	s, ok := value.(string)
	if !ok {
		data, _ := json.Marshal(value)
		return string(data)
	}

	return s
}

// CallerFormat format entry caller as text with fmt.Sprintf
func CallerFormat(caller *runtime.Frame) string {
	i := strings.LastIndex(caller.File, "/")
	i = strings.LastIndex(caller.File[:i], "/")
	return fmt.Sprintf("%s:%d", caller.File[i+1:], caller.Line)
}

// NeedsQuoting checks quoting needed for text
func NeedsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
