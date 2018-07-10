package echo

import (
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// Skipper defines a function to skip middleware.
type Skipper func(c echo.Context) bool

// Config defines the config for Logger middleware.
type Config struct {
	Logger *logrus.Logger
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Fields available for logging
	// - id (Request ID)
	// - ip
	// - host
	// - referer
	// - user_agent
	// - status
	// - latency
	Fields []string
}

var (
	// DefaultConfig is the default Logger middleware config.
	DefaultConfig = Config{
		Logger:  logrus.StandardLogger(),
		Skipper: func(c echo.Context) bool { return false },
		Fields:  []string{"ip", "latency", "status"},
	}
)

// Middleware returns a Logger middleware with config.
func Middleware(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.Logger == nil {
		config.Logger = logrus.StandardLogger()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}
			fields := logrus.Fields{}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
				fields["error"] = err
			}
			stop := time.Now()

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			for _, field := range config.Fields {
				switch field {
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					fields[field] = id
				case "ip":
					fields[field] = c.RealIP()
				case "host":
					fields[field] = req.Host
				case "referer":
					fields[field] = req.Referer()
				case "user_agent":
					fields[field] = req.UserAgent()
				case "status":
					fields[field] = res.Status
				case "latency":
					fields[field] = stop.Sub(start).String()
				}
			}

			switch {
			case res.Status >= 500:
				config.Logger.WithFields(fields).Errorf("%s %s", req.Method, path)
			case res.Status >= 400:
				config.Logger.WithFields(fields).Warnf("%s %s", req.Method, path)
			default:
				config.Logger.WithFields(fields).Infof("%s %s", req.Method, path)
			}

			return err
		}
	}
}
