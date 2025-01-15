package api

import (
	"context"
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiCallTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCallTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCallTemplateLogic {
	return &ApiCallTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCallTemplateLogic) ApiCallTemplate(req *types.ApiCallTemplateRequest) (resp *types.ApiCallTemplateResponse, err error) {
	// 查询 api 参数
	api, err := l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, err
	}
	url := l.svcCtx.Config.ApiUrl + "/" + req.ApiId
	// 查询 canvas
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, api.WorkspaceId)
	if err != nil {
		return nil, err
	}
	param, err := l.readData(gjson.Parse(canvas.Draft))
	if err != nil {
		return nil, err
	}
	// 查询 secret
	header := map[string]string{}
	secret, err := l.svcCtx.ApiSecretKeyModel.FindByApiId(l.ctx, api.ApiId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if secret == nil {
		header["Authorization"] = "Bearer " + "请创建密钥"
	} else {
		header["Authorization"] = "Bearer " + secret[0].SecretKey
	}
	headerJson, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	resp = &types.ApiCallTemplateResponse{
		ApiId:  api.ApiId,
		Url:    url,
		Header: string(headerJson),
		Body:   param,
	}

	return resp, nil
}

func (l *ApiCallTemplateLogic) readData(result gjson.Result) (string, error) {
	nodes := result.Get("graph.nodes").Array()
	for _, node := range nodes {
		if node.Get("data.type").String() == "start" {
			param := node.Get("data.custom.param").String()
			var data interface{}
			if err := json.Unmarshal([]byte(param), &data); err == nil {
				return param, nil
			}
			return "", errors.New(int(logic.SystemError), "输入不是 JSON 格式")
		}
	}
	return "", errors.New(int(logic.SystemError), "未找到开始节点")
}
