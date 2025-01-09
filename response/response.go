package response

import (
	errors2 "errors"
	"github.com/zeromicro/x/errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		var codeMsg *errors.CodeMsg
		if errors2.As(err, &codeMsg) {
			body.Code = codeMsg.Code
			body.Msg = codeMsg.Msg
		} else {
			body.Code = -1
			body.Msg = err.Error()
		}
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
