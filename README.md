# Hooks for [logrus](https://github.com/Sirupsen/logrus)

Example
```go
package main

import (
  "github.com/onrik/logrus/filename"
  log "github.com/Sirupsen/logrus"
)

func main() {
  log.AddHook(filename.NewHook(log.ErrorLevel, log.FatalLevel))
}

```
