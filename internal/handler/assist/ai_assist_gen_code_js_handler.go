package assist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"workflow/internal/logic/assist"
	"workflow/internal/svc"
	"workflow/internal/types"
)

func AiAssistGenCodeJsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 SSE 必需的 HTTP 头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		var req types.AiAssistGenCodeJsRequest
		if err := httpx.Parse(r, &req); err != nil {
			sendError(w, err)
			return
		}

		// 为每个客户端创建一个 channel
		clientChan := make(chan string)
		errorChan := make(chan error)
		doneChan := make(chan struct{})

		// 客户端断开时清理
		defer func() {
			close(clientChan)
			close(errorChan)
			close(doneChan)
		}()

		l := assist.NewAiAssistGenCodeJsLogic(r.Context(), svcCtx)
		go func() {
			err := l.AiAssistGenCodeJs(&req, clientChan)
			if err != nil {
				errorChan <- err
				return
			}
			doneChan <- struct{}{}
		}()

		// 持续监听并推送事件
		for {
			select {
			case msg := <-clientChan:
				jsonData, _ := json.Marshal(map[string]string{"content": msg})
				// 发送数据事件
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
				w.(http.Flusher).Flush()
			case err := <-errorChan:
				sendError(w, err)
				return
			case <-doneChan:
				// 发送完成事件
				fmt.Fprintf(w, "event: done\ndata: {}\n\n")
				w.(http.Flusher).Flush()
				return
			case <-r.Context().Done():
				// 客户端断开连接
				return
			}
		}
	}
}

func sendError(w http.ResponseWriter, err error) {
	errorData := map[string]string{"error": err.Error()}
	jsonData, _ := json.Marshal(errorData)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", jsonData)
	w.(http.Flusher).Flush()
}
