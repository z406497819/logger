# logger
my file logger system

```go
import (
	"tinky.golang.com/mongo_driver/conf"
)

var log logger.Logger

log = logger.NewFileLogger("debug", "./log/")

log.Info("hahahaha")
log.Debug("hahahaha")
log.Warning("hahahaha")
log.Error("hahahaha")
```