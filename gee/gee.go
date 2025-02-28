package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 定义了框架使用的请求处理程序，遵循标准库的HandlerFunc格式
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 是框架的核心结构，管理路由表和处理函数映射
// router 是一个映射表，键为由HTTP方法和路径组成的字符串（如"GET-/hello"）
// 值为对应的处理函数HandlerFunc
type Engine struct {
	router map[string]HandlerFunc
}

// New 创建并返回一个新的Engine实例，初始化路由映射表
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute 私有方法，用于向路由表中添加新的路由规则
// method: HTTP方法（GET/POST等）
// pattern: 请求路径
// handler: 对应的处理函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET 注册GET方法的路由处理函数
// 将GET请求路径与处理函数绑定到路由表
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 注册POST方法的路由处理函数
// 将POST请求路径与处理函数绑定到路由表
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 启动HTTP服务器并监听指定地址
// 使用Engine实例作为请求处理器
// 返回可能的服务器错误
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 实现http.Handler接口，处理所有HTTP请求
// 根据请求方法和路径生成路由键，查找对应的处理函数
// 找到则执行处理函数，否则返回404状态和错误信息
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		// 未找到路由时返回404
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
