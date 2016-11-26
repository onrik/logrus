# Hooks for [logrus](https://github.com/Sirupsen/logrus)

Example
```go
package main

import (
  "github.com/onrik/logrus/filename"
  "github.com/onrik/logrus/sentry"
  log "github.com/Sirupsen/logrus"
)

func main() {
  dsn := "http://60a0257d7b5a429a8838e5f2ba873ec9:cb785a64cd3649ea987a1f2f5fad5e82@example.com/2"
  
  log.AddHook(filename.NewHook(log.ErrorLevel, log.FatalLevel))
  log.AddHook(sentry.NewHook(dsn))
}

```
