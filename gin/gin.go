package gin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger *logrus.Logger

	// Fields available for logging
	// - id (Request ID)
	// - ip
	// - host
	// - referer
	// - user_agent
	// - status
	// - latency
	// - headers
	Fields []string
}

var (
	DefaultConfig = Config{
		Logger: logrus.StandardLogger(),
		Fields: []string{"ip", "latency", "status"},
	}
)

func Middleware(config Config) gin.HandlerFunc {
	if config.Logger == nil {
		config.Logger = DefaultConfig.Logger
	}
	if config.Fields == nil {
		config.Fields = DefaultConfig.Fields
	}

	return func(c *gin.Context) {
		start := time.Now()

		request := c.Request
		method := request.Method
		path := request.URL.Path
		if path == "" {
			path = "/"
		}
		c.Next()

		stop := time.Now()

		fields := logrus.Fields{}
		status := c.Writer.Status()
		for _, field := range config.Fields {
			switch field {
			case "ip":
				fields[field] = c.ClientIP()
			case "host":
				fields[field] = request.Host
			case "referer":
				fields[field] = request.Referer()
			case "user_agent":
				fields[field] = request.UserAgent()
			case "status":
				fields[field] = status
			case "latency":
				fields[field] = stop.Sub(start).String()
			case "headers":
				fields[field] = request.Header
			}
		}

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		switch {
		case status >= 500:
			config.Logger.WithFields(fields).Errorf("%s %s", method, path)
		case status >= 400:
			config.Logger.WithFields(fields).Warnf("%s %s", method, path)
		default:
			config.Logger.WithFields(fields).Debugf("%s %s", method, path)
		}
	}
}
