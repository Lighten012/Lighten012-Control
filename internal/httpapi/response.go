// Package httpapi 提供 HTTP API 层：统一 JSON 响应格式、
// 路由组装（挂载各功能模块）与静态页托管。
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
)

// PathPrefix 所有功能模块 API 路由的统一根路径。
const PathPrefix = "/lighten012-api"

// Response 统一 API 响应结构。
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// WriteOK 输出成功响应（code=0）。
func WriteOK(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// WriteError 输出错误响应。
func WriteError(w http.ResponseWriter, httpStatus int, message string) {
	write(w, httpStatus, Response{Code: httpStatus, Message: message})
}

// RegisterGet 注册一个只读查询路由：执行查询并把结果按统一 JSON 格式写出。
func RegisterGet[T any](mux *http.ServeMux, path string, fn func(ctx context.Context) (T, error)) {
	mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
		data, err := fn(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		WriteOK(w, data)
	})
}

func write(w http.ResponseWriter, httpStatus int, body Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}
