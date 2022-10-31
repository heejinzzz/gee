# gee
* gee 是一个模仿 [gin](https://github.com/gin-gonic/gin) 编写的 web 框架
* 支持路由注册、路由分组、中间件、cookie、html渲染、日志记录、错误恢复等

---
### Basic Usage
gee 的语法与 gin 基本一致，下面是一个基本的使用示例：
```go
package main

import (
	"github.com/heejinzzz/gee"
	"net/http"
)

func main() {
	r := gee.NewDefault()
	
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello, Gee!")
	})

	r.Run(":8080")
}
```