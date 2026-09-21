// Package httpapi 提供 HTTP API 层：统一 JSON 响应格式、
// 路由组装（挂载各功能模块）与静态页托管。
package httpapi

import (
	"encoding/json"
	"net/http"
)

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

func write(w http.ResponseWriter, httpStatus int, body Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}
