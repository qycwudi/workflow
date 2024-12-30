package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// 自定义的中间件
func AuditMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求 url 和 body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 重新设置 body
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		logx.Infof("request url: %s, body: %s", r.URL, string(body))

		// 创建一个自定义的 ResponseWriter 来捕获响应
		rw := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next(rw, r)

		// 读取响应内容和状态码
		respBody := rw.body.Bytes()
		logx.Infof("response status: %d, body: %s", rw.statusCode, string(respBody))
	}
}

// 自定义的 ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// Write 方法同时写入到原始 ResponseWriter 和缓冲区
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader 方法记录状态码
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
