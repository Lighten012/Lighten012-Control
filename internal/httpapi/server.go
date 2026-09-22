package httpapi

import (
	"log"
	"net/http"
	"time"
)

// MountAPI 为未匹配的 API 路径补充统一 JSON 404。
func MountAPI(mux *http.ServeMux) {
	mux.HandleFunc(PathPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotFound, "接口不存在")
	})
}

// MountStatic 托管内置静态页面。
func MountStatic(mux *http.ServeMux, static http.FileSystem) {
	if static != nil {
		mux.Handle("/", http.FileServer(static))
	}
	log.Printf("[server] 已挂载静态页面")
}

// WithLogging 输出简要访问日志（方法、路径、状态码、耗时）。
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		log.Printf("[http] %s %s -> %d (%s)",
			r.Method, r.URL.Path, sw.status, time.Since(start).Round(time.Millisecond))
	})
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
