package logic

type Code int

const (
	SUCCESS                         Code = 0
	KeyExist                        Code = 100
	KeyMiss                         Code = 101
	SendMessageErr                  Code = 102
	SendMessageParamFormattingError Code = 103

	// SystemError 系统异常
	SystemError Code = 500
	// SystemStoreError 系统存储异常
	SystemStoreError Code = 501
	// SystemOrmError 系统Orm异常
	SystemOrmError Code = 502
	// ParamError 参数异常
	ParamError Code = 503
)
