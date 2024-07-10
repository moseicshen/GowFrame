package gow

import (
	"fmt"
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		fmt.Printf("Path: %s\n", c.Path)
		fmt.Printf("Method: %s\n", c.Method)
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
