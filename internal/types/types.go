// Code generated by goctl. DO NOT EDIT.
package types

type SetKvRequest struct {
	SpiderName string `json:"spiderName"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

type SetKvResponse struct {
	Code    int    `json:"code" common:"0-成功|100-key存在|500-系统错误"`
	Message string `json:"message"`
}

type GetVByKRequest struct {
	Key string `json:"key"`
}

type GetVByKResponse struct {
	Code    int    `json:"code" common:"0-成功|101-key不存在|500-系统错误"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WorkSpaceRequest struct {
	Id            string `json:"id,omitempty"`
	WorkSpaceName string `json:"workSpaceName"`
}

type WorkSpaceNewResponse struct {
	Response
	Id string `json:"id"`
}

type WorkSpaceEditResponse struct {
	Response
	Id string `json:"id"`
}

type WorkRemoveRequest struct {
	Id string `json:"id"`
}

type WorkSpaceRemoveResponse struct {
	Response
	Result bool `json:"result"`
}
