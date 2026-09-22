package httpapi

import (
	"log"
	"net/http"
	"time"
)

// PathPrefix 是所有功能 API 的统一根路径。
const PathPrefix = "/lighten012-api"

// New 组装根 HTTP 处理器：先由入口代码挂载具体功能路由，
// 未匹配的 API 路径统一返回 JSON 404，其余路径由静态页托管。
func New(static http.FileSystem, registerRoutes func(mux *http.ServeMux)) http.Handler {
	mux := http.NewServeMux()
	registerRoutes(mux)

	mux.HandleFunc(PathPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotFound, "接口不存在")
	})

	if static != nil {
		mux.Handle("/", http.FileServer(static))
	}

	return withLogging(mux)
}

// statusWriter 记录响应状态码，用于访问日志。
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// withLogging 输出简要访问日志（方法、路径、状态码、耗时）。
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		log.Printf("[http] %s %s -> %d (%s)",
			r.Method, r.URL.Path, sw.status, time.Since(start).Round(time.Millisecond))
	})
}
