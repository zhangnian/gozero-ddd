package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware 日志中间件
// 演示 go-zero 中间件的实现方式
type LoggingMiddleware struct{}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Handle 处理请求
func (m *LoggingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 调用下一个处理器
		next(w, r)

		// 记录请求日志
		log.Printf(
			"[%s] %s %s - %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	}
}

