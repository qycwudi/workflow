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
)
