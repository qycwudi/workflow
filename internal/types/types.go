// Code generated by goctl. DO NOT EDIT.
package types

type WorkSpaceBase struct {
	Id            string   `json:"id,optional"`
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
	Id string `json:"id,optional"`
}

type WorkSpaceRemoveResponse struct {
}

type WorkSpaceListRequest struct {
	WorkSpaceName string  `json:"workSpaceName,optional"`
	WorkSpaceType string  `json:"workSpaceType,optional"`
	WorkSpaceTag  []int64 `json:"workSpaceTag,optional"`
	Current       int     `json:"current"`
	PageSize      int     `json:"pageSize"`
}

type WorkSpaceListResponse struct {
	Current  int             `json:"current"`
	PageSize int             `json:"pageSize"`
	Total    int64           `json:"total"`
	Data     []WorkSpacePage `json:"data"`
}

type WorkSpacePage struct {
	WorkSpaceBase
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type WorkSpaceEditTagRequest struct {
	Id           string   `json:"id"`
	WorkSpaceTag []string `json:"workSpaceTag"`
}

type WorkSpaceEditTagResponse struct {
}

type TagListRequest struct {
}

type TagListResponse struct {
	Tag []TagEntity `json:"tagList"`
}

type TagEntity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MockRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MockInfo struct {
	Address   string   `json:"address"`
	PhoneNums []string `json:"phoneNums"`
}

type MockResponse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CanvasRunRequest struct {
	Id string `json:"id" desc:"空间ID"`
}

type CanvasRunResponse struct {
	Ts       int64                  `json:"ts"`
	Id       string                 `json:"id"`
	MetaData map[string]string      `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}

type CanvasRunSingleRequest struct {
	Id     string `json:"id" desc:"空间ID"`
	NodeId string `json:"nodeId" desc:"节点ID"`
}

type CanvasRunSingleResponse struct {
	Ts       int64                  `json:"ts"`
	Id       string                 `json:"id"`
	MetaData map[string]string      `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}

type GetCanvasRunHistoryReq struct {
	WorkSpaceId string `path:"workSpaceId"`
}

type RunHistoryRecord struct {
	Id             string `json:"id"`             // 运行记录ID
	StartTime      string `json:"startTime"`      // 开始时间，ISO格式
	Duration       int64  `json:"duration"`       // 总耗时(ms)
	Status         string `json:"status"`         // 运行状态
	ComponentCount int64  `json:"componentCount"` // 组件数量
}

type GetCanvasRunHistoryResp struct {
	Records []RunHistoryRecord `json:"records"`
	Total   int64              `json:"total"` // 总记录数
}

type GetCanvasRunDetailReq struct {
	RecordId string `path:"recordId"` // 运行记录ID
}

type GetCanvasRunDetailResp struct {
	Id         string            `json:"id"`         // 运行记录ID
	StartTime  string            `json:"startTime"`  // 开始时间
	Duration   int64             `json:"duration"`   // 总耗时(ms)
	Status     string            `json:"status"`     // 运行状态 success/failed
	Error      string            `json:"error"`      // 错误信息
	Components []ComponentDetail `json:"components"` // 组件列表
}

type ComponentDetail struct {
	Id        string                 `json:"id"`        // 组件ID
	Name      string                 `json:"name"`      // 组件名称
	Logic     string                 `json:"logic"`     // 组件类型
	StartTime int64                  `json:"startTime"` // 开始时间戳
	Duration  int64                  `json:"duration"`  // 耗时(ms)
	Status    string                 `json:"status"`    // 组件运行状态 success/failed
	Error     string                 `json:"error"`     // 组件错误信息
	Input     map[string]interface{} `json:"input"`     // 输入参数
	Output    map[string]interface{} `json:"output"`    // 输出结果
}

type CanvasDraftRequest struct {
	Id    string                 `json:"id"`
	Graph map[string]interface{} `json:"graph"`
}

type CanvasDraftResponse struct {
	Hash       string `json:"hash"`
	UpdateTime int64  `json:"updateTime"`
}

type EdgeCustomData struct {
	SourcePoint int `json:"sourcePoint"`
	TargetPoint int `json:"targetPoint"`
}

type CanvasPublishRequest struct {
	Id      string `json:"id" desc:"空间ID"`
	ApiName string `json:"apiName" desc:"名称"`
	ApiDesc string `json:"apiDesc" desc:"描述"`
}

type CanvasPublishResponse struct {
	ApiId string `json:"apiId"`
}

type CanvasDetailRequest struct {
	Id string `json:"id"`
}

type CanvasDetailResponse struct {
	Id    string                 `json:"id"`
	Graph map[string]interface{} `json:"graph"`
}

type ApiPublishListRequest struct {
	Current  int    `json:"current"`
	PageSize int    `json:"pageSize"`
	Id       string `json:"id,optional" desc:"空间ID 非必填"`
}

type ApiPublishListResponse struct {
	Current  int              `json:"current"`
	PageSize int              `json:"pageSize"`
	Total    int64            `json:"total"`
	List     []ApiPublishList `json:"list"`
}

type ApiPublishList struct {
	WorkSpaceId string `json:"workSpaceId"`
	ApiId       string `json:"apiId"`
	ApiName     string `json:"apiName"`
	ApiDesc     string `json:"apiDesc"`
	Status      string `json:"status" desc:"上下线状态ON OFF"`
}

type ApiOnOffRequest struct {
	ApiId  string `json:"apiId"`
	Status string `json:"status" desc:"上下线状态ON OFF"`
}

type ApiOnOffResponse struct {
	ApiId  string `json:"apiId"`
	Status string `json:"status" desc:"上下线状态ON OFF"`
}

type ApiRecordsRequest struct {
	Current  int    `json:"current"`
	PageSize int    `json:"pageSize"`
	ApiId    string `json:"apiId,optional" desc:"apiId"`
	ApiName  string `json:"apiName,optional" desc:"api名称"`
}

type ApiRecordsResponse struct {
	Current  int          `json:"current"`
	PageSize int          `json:"pageSize"`
	Total    int64        `json:"total"`
	List     []ApiRecords `json:"list"`
}

type ApiRecords struct {
	ApiId    string `json:"apiId"`
	ApiName  string `json:"apiName" desc:"名称"`
	CallTime string `json:"callTime"`
	Status   string `json:"status"`
	TraceId  string `json:"traceId"`
	Param    string `json:"param"`
	Extend   string `json:"extend"`
}

type ApiSecretyKeyListRequest struct {
	ApiId string `json:"apiId"`
}

type ApiSecretyKeyListResponse struct {
	Total int64           `json:"total"`
	List  []ApiSecretyKey `json:"list"`
}

type ApiSecretyKey struct {
	ApiId          string `json:"apiId"`
	SecretyKey     string `json:"secretyKey"`
	ExpirationTime string `json:"expirationTime"`
}

type ModuleListRequest struct {
}

type ModuleListResponse struct {
	Total   int          `json:"total"`
	Modules []ModuleData `json:"modules"`
}

type ModuleData struct {
	Index        int    `json:"index"`
	ModuleId     string `json:"moduleId"`
	ModuleName   string `json:"moduleName"`
	ModuleType   string `json:"moduleType"`
	ModuleConfig string `json:"moduleConfig"`
}

type ModuleNewRequest struct {
	Index        int    `json:"index"`
	ModuleName   string `json:"moduleName"`
	ModuleType   string `json:"moduleType"`
	ModuleConfig string `json:"moduleConfig"`
}

type ModuleNewResponse struct {
	ModuleId string `json:"moduleId"`
}

type ModuleEditRequest struct {
	Index        int    `json:"index"`
	ModuleId     string `json:"moduleId"`
	ModuleName   string `json:"moduleName"`
	ModuleType   string `json:"moduleType"`
	ModuleConfig string `json:"moduleConfig"`
}

type ModuleEditResponse struct {
	ModuleId string `json:"moduleId"`
}

type DatasourceListRequest struct {
	Current  int    `json:"current"`
	PageSize int    `json:"pageSize"`
	Type     string `json:"type,optional"`
	Status   string `json:"status,optional"`
	Switch   int    `json:"switch,optional"`
}

type DatasourceListResponse struct {
	Current  int              `json:"current"`
	PageSize int              `json:"pageSize"`
	Total    int64            `json:"total"`
	List     []DatasourceInfo `json:"list"`
}

type DatasourceInfo struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	Config     string `json:"config"`
	Switch     int    `json:"switch"`
	Hash       string `json:"hash"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type DatasourceAddRequest struct {
	Type   string `json:"type"`
	Config string `json:"config"`
	Switch int    `json:"switch"`
	Hash   string `json:"hash,optional"`
	Status string `json:"status,optional"`
}

type DatasourceAddResponse struct {
	Id int `json:"id"`
}

type DatasourceEditRequest struct {
	Id     int    `json:"id"`
	Type   string `json:"type,optional"`
	Config string `json:"config,optional"`
	Switch int    `json:"switch,optional"`
	Hash   string `json:"hash,optional"`
	Status string `json:"status,optional"`
}

type DatasourceEditResponse struct {
	Id int `json:"id"`
}

type DatasourceDeleteRequest struct {
	Id int `json:"id"`
}

type DatasourceDeleteResponse struct {
	Id int `json:"id"`
}

type DatasourceTestRequest struct {
	Type   string `json:"type,optional"`
	Config string `json:"config,optional"`
	Switch int    `json:"switch,optional"`
}

type DatasourceTestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
