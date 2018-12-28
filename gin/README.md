# Logger for [gin](https://github.com/gin-gonic/gin)

Example
```go
package main

import (
    "github.com/gin-gonic/gin"
    ginlog "github.com/onrik/logrus/gin"
    "github.com/sirupsen/logrus"
)

func main() {
    server := gin.New()

    server.Use(ginlog.Middleware(ginlog.DefaultConfig))
}
```
