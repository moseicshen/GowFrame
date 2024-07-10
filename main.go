package main

import (
	"gow"
	"log"
	"net/http"
	"time"
)

func forV2() gow.HandlerFunc {
	return func(c *gow.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("v2-[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gow.New()
	r.Use(gow.Logger())
	r.GET("/", func(c *gow.Context) {
		c.JSON(200, gow.H{
			"message": "Hello Gow",
		})
	})
	v2 := r.Group("/v2")
	v2.Use(forV2())
	v2.GET("/hello", func(c *gow.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.ParamValue("name"), c.Path)
	})
	r.Run(":9999")
}
