# Logger for [gorm](http://doc.gorm.io/)

Example
```go
package main

import (
    "github.com/jinzhu/gorm"
    gormlog "github.com/onrik/logrus/gorm"
    "github.com/sirupsen/logrus"
)

func main() {
    db, err := gorm.Open("<driver>", "<dsn>")
    if err != nil {
        logrus.Fatal(err)
    }

    db.SetLogger(gormlog.New(logrus.StandardLogger()))
    db.LogMode(true) // Will log only on logrus debug level
}
```
