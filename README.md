# logger
my file logger system.

```go
package main

import (
	"github.com/z406497819/logger"
)

var log logger.Logger

func main() {
	log = logger.NewFileLogger("debug", "./log/", false)
	
	log.Info("log message:%s","hello world")
}


```