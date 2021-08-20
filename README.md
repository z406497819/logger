# logger
my file logger system.

```go
package main

import (
	"github.com/z406497819/logger"
)

func main() {

    // 如果需要额外配置，使用如下操作即可
    // logger.AddOption(logger.WithLevelStr("debug"), logger.WithPath("./log1"), logger.WithAsync(true))

    logger.Info("log message:%s", "hello world")
    logger.Error("log message:%s", "hello world")
}


```
