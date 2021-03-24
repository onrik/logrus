# Hooks for [logrus](https://github.com/Sirupsen/logrus)

Example
```go
package main

import (
    "fmt"

    "github.com/onrik/logrus/filename"
    "github.com/onrik/logrus/sentry"
    log "github.com/sirupsen/logrus"
)

var (
    dsn = "http://60a0257d7b5a429a8838e5f2ba873ec9@example.com/1"
)

func main() {
    filenameHook := filename.NewHook()
    filenameHook.Field = "custom_source_field" // Customize source field name
    log.AddHook(filenameHook)

    sentryHook, err := sentry.NewHook(sentry.Options{
        Dsn: dsn,
    })
    if err != nil {
        log.Error(err)
        return
    }
    log.AddHook(sentryHook, log.PanicLevel, log.FatalLevel, log.ErrorLevel)

    err = fmt.Errorf("test error")
    log.WithError(err).Error("Dead beef")
}


```
