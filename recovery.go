package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var builder strings.Builder
	builder.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		builder.WriteString(fmt.Sprintf("\n\t%s: Line %d", file, line))
	}
	return builder.String()
}

func Recovery(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Printf("%s\n\n", trace(message))
			c.Fail(http.StatusInternalServerError, "Internal Server Error")
		}
	}()

	c.Next()
}
