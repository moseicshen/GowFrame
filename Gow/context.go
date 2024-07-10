package gow

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// original objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Param  map[string]string
	// response info
	StatusCode int
	// middleware
	handlers  []HandlerFunc
	index     int
	midStatus bool
	// engine
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:    w,
		Req:       req,
		Path:      req.URL.Path,
		Method:    req.Method,
		index:     -1,
		midStatus: true,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
		if !c.midStatus {
			break
		}
	}
}

func (c *Context) Abort() {
	c.midStatus = false
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) ParamValue(key string) string {
	value, _ := c.Param[key]
	return value
}

// PostForm 获取表单变量的参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取URL中的请求参数（包括加密项）
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// JSON 编码对象，向客户端返回json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 创建json编码器
	encoder := json.NewEncoder(c.Writer)
	// 通过编码器将结构体编码为json，传输至客户端
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

// Data 返回数据流
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		return
	}
}

// String 返回字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		return
	}
}

// HTML 返回HTML格式数据
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		return
	}
}
