package main

import (
	"gow"
	"net/http"
)

func main() {
	r := gow.New()
	r.GET("/index", func(c *gow.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gow.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gow</h1>")
		})

		v1.GET("/hello", func(c *gow.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello v1 %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gow.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello v2 %s, you're at %s\n", c.ParamValue("name"), c.Path)
		})
		v2.POST("/login", func(c *gow.Context) {
			c.JSON(http.StatusOK, gow.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
