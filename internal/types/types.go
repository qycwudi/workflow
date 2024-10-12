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

type WorkSpaceDetailRequest struct {
	Id string `json:"id"`
}

type WorkSpaceDetailResponse struct {
	Id                    string                 `json:"id"`
	Graph                 map[string]interface{} `json:"graph"`
	Features              map[string]interface{} `json:"features"`
	EnvironmentVariables  []interface{}          `json:"environment_variables"`
	ConversationVariables []interface{}          `json:"conversation_variables"`
	Hash                  string                 `json:"hash"`
	BaseInfo              WorkSpaceBase          `json:"baseInfo"`
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
	Id       string                 `json:"id" desc:"空间ID"`
	MetaData map[string]string      `json:"metaData" desc:"元数据"`
	Data     map[string]interface{} `json:"data" desc:"数据"`
}

type CanvasRunResponse struct {
	Ts       int64                  `json:"ts"`
	Id       string                 `json:"id"`
	MetaData map[string]string      `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}

type CanvasRunRecordRequest struct {
	Id string `json:"id" desc:"空间ID"`
}

type CanvasRunRecordResponse struct {
	Records []RunRecord `json:"records"`
}

type RunRecord struct {
	TraceId string `json:"traceId"`
	Status  string `json:"status"`
	RunTime string `json:"runTime" desc:"开始时间"`
}

type TraceRequest struct {
	TraceId string `json:"traceId"`
}

type TraceResponse struct {
	Total            int64   `json:"total"`
	TotalElapsedTime int64   `json:"total"`
	Traces           []Trace `json:"traces"`
}

type Trace struct {
	TraceId     string `json:"traceId"`
	NodeId      string `json:"nodeId" desc:"节点ID"`
	NodeName    string `json:"nodeName" desc:"节点名称"`
	Status      string `json:"status" desc:"运行状态"`
	StartTime   string `json:"startTime" desc:"开始执行时间"`
	ElapsedTime int64  `json:"elapsedTime" desc:"运行耗时"`
	Input       string `json:"input" desc:"输出"`
	Output      string `json:"output" desc:"输出"`
	Logic       string `json:"logic" desc:"执行逻辑"`
}

type CanvasDraftRequest struct {
	Id                    string                 `json:"id"`
	Graph                 map[string]interface{} `json:"graph"`
	Features              map[string]interface{} `json:"features"`
	EnvironmentVariables  []interface{}          `json:"environment_variables"`
	ConversationVariables []interface{}          `json:"conversation_variables"`
	Hash                  string                 `json:"hash,optional"`
}

type CanvasDraftResponse struct {
	Hash       string `json:"hash"`
	UpdateTime int64  `json:"updateTime"`
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

type EdgeCustomData struct {
	SourcePoint int `json:"sourcePoint"`
	TargetPoint int `json:"targetPoint"`
}

type CanvasDetailRequest struct {
	Id string `json:"id"`
}

type CanvasDetailResponse struct {
	Id                    string                 `json:"id"`
	Graph                 map[string]interface{} `json:"graph"`
	Features              map[string]interface{} `json:"features"`
	EnvironmentVariables  []interface{}          `json:"environment_variables"`
	ConversationVariables []interface{}          `json:"conversation_variables"`
	Hash                  string                 `json:"hash"`
}
