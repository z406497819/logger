# logger
my file logger system.

```go
package main

import (
	"github.com/z406497819/logger"
)

func main() {

    // 如果需要额外配置，使用如下操作即可
    // AddOption(WithLevelStr("debug"), WithPath("./log1"), WithAsync(true))
    
    Info("log message:%s", "hello world")
}


```