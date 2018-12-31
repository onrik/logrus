# Logger for [echo](https://echo.labstack.com/)

Example
```go
package main

import (
    "github.com/labstack/echo"
    echolog "github.com/onrik/logrus/echo"
    "github.com/sirupsen/logrus"
)

func main() {
    server := echo.New()

    server.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
    server.Use(echolog.Middleware(echolog.DefaultConfig))
}
```
