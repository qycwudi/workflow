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

type WorkSpaceBase struct {
	WorkSpaceId   string   `json:"workSpaceId,optional"`
	WorkSpaceName string   `json:"workSpaceName"`
	WorkSpaceDesc string   `json:"workSpaceDesc,optional"`
	WorkSpaceType string   `json:"workSpaceType"`
	WorkSpaceTag  []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string   `json:"workSpaceIcon,optional"`
}

type WorkSpaceNewRequest struct {
	WorkSpaceBase
}

type WorkSpaceNewResponse struct {
	WorkSpaceBase
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkSpaceEditRequest struct {
	WorkSpaceBase
	WorkSpaceConfig string `json:"workSpaceConfig,optional"`
}

type WorkSpaceEditResponse struct {
	WorkSpaceBase
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkRemoveRequest struct {
	WorkSpaceId string `json:"workSpaceId,optional"`
}

type WorkSpaceRemoveResponse struct {
}

type WorkSpaceListRequest struct {
	WorkSpaceName string   `json:"workSpaceName,optional"`
	WorkSpaceType string   `json:"workSpaceType,optional"`
	WorkSpaceTag  []string `json:"workSpaceTag,optional"`
	Current       int      `json:"current"`
	PageSize      int      `json:"pageSize"`
}

type WorkSpaceListResponse struct {
	Current  int             `json:"current"`
	PageSize int             `json:"pageSize"`
	Total    int             `json:"total"`
	Data     []WorkSpacePage `json:"data"`
}

type WorkSpacePage struct {
	WorkSpaceBase
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type WorkSpaceDetailRequest struct {
	WorkSpaceId string `json:"workSpaceId"`
}

type WorkSpaceDetailResponse struct {
	WorkSpaceBase
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkSpaceEditTagRequest struct {
	WorkSpaceId  string   `json:"workSpaceId"`
	WorkSpaceTag []string `json:"workSpaceTag"`
}

type WorkSpaceEditTagResponse struct {
}

type WorkSpaceUploadCanvasConfigTagRequest struct {
	WorkSpaceId  string `json:"workSpaceId"`
	CanvasConfig string `json:"canvasConfig"`
}

type WorkSpaceUploadCanvasConfigTagResponse struct {
}
