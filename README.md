# logger
my file logger system

```go
import (
	"github.com/z406497819/logger"
)

var log logger.Logger

func main()
{
    log = logger.NewFileLogger("debug", "./log/")
    
    log.Info("hahahaha")
    log.Debug("hahahaha")
    log.Warning("hahahaha")
    log.Error("hahahaha")
}

```