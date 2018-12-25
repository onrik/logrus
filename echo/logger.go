package echo

import (
	"io"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var levelsMap = map[log.Lvl]logrus.Level{
	log.DEBUG: logrus.DebugLevel,
	log.INFO:  logrus.InfoLevel,
	log.WARN:  logrus.WarnLevel,
	log.ERROR: logrus.ErrorLevel,
	log.OFF:   logrus.FatalLevel,
}

type Logger struct {
	*logrus.Logger
	prefix   string
	level    log.Lvl
	MsgField string
}

// NewLogger creates logger
func NewLogger(l *logrus.Logger, prefix string) *Logger {
	logger := &Logger{
		Logger:   l,
		prefix:   prefix,
		MsgField: "",
	}

	switch l.GetLevel() {
	case logrus.DebugLevel, logrus.TraceLevel:
		logger.level = log.DEBUG
	case logrus.InfoLevel:
		logger.level = log.INFO
	case logrus.WarnLevel:
		logger.level = log.WARN
	case logrus.ErrorLevel:
		logger.level = log.ERROR
	case logrus.FatalLevel, logrus.PanicLevel:
		logger.level = log.OFF
	default:
		logger.level = log.INFO
	}

	return logger
}

func (l *Logger) fields(j log.JSON) (string, logrus.Fields) {
	msg, _ := j[l.MsgField].(string)
	delete(j, l.MsgField)
	return l.Prefix() + msg, logrus.Fields(j)
}

func (l *Logger) SetHeader(h string) {
}

func (l *Logger) Output() io.Writer {
	return l.Out
}

func (l *Logger) SetOutput(w io.Writer) {
	l.SetOutput(w)
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Level() log.Lvl {
	return l.level
}

func (l *Logger) SetLevel(lvl log.Lvl) {
	l.level = lvl
	l.Logger.SetLevel(levelsMap[lvl])
}

func (l *Logger) Printj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Print(msg)
}

func (l *Logger) Debugj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Debug(msg)
}

func (l *Logger) Infoj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Info(msg)
}

func (l *Logger) Warnj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Warn(msg)
}

func (l *Logger) Errorj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Error(msg)
}

func (l *Logger) Fatalj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Fatal(msg)
}

func (l *Logger) Panicj(j log.JSON) {
	msg, fields := l.fields(j)
	l.WithFields(fields).Panic(msg)
}
