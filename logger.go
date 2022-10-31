package gee

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const MaxInvisibleLengthOfURI = 30

func Logger(c *Context) {
	// start time
	t := time.Now()
	// process all middlewares and handler
	c.Next()
	// calculate process duration
	duration := time.Since(t)

	method := fmt.Sprintf("[%s]", c.Method)
	path := c.Req.RequestURI
	if len(path) > MaxInvisibleLengthOfURI {
		path = path[:MaxInvisibleLengthOfURI] + "..."
	}
	statusCodeLen := len(strconv.Itoa(c.StatusCode))
	statusCodeLeftSpace := (9 - statusCodeLen) / 2
	statusCodeRightSpace := 9 - statusCodeLeftSpace - statusCodeLen

	format := " %-" + strconv.Itoa(statusCodeLeftSpace) + "s%d%" + strconv.Itoa(statusCodeRightSpace) + "s  %-6s  %-40s %v\n"
	log.Printf(format, "|", c.StatusCode, "|", method, path, duration)
}
