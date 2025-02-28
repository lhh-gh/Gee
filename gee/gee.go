package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 定义请求处理函数类型，接收上下文对象
type HandlerFunc func(*Context)

// RouterGroup 路由组结构，支持路由分组和中间件
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// Engine 核心结构，实现http.Handler接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// New 创建并初始化引擎实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 创建子路由组（支持路由分组嵌套）
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute 注册路由到路由树
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 注册GET方法路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 注册POST方法路由（
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 为路由组添加中间件
// authGroup.Use(JWTAuth(), Logging())

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP 实现http.Handler接口（核心请求处理方法）
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 收集匹配请求路径的所有中间件
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 创建请求上下文对象
	c := newContext(w, req)
	c.handlers = middlewares // 注入中间件链
	engine.router.handle(c)  // 执行路由处理
}
