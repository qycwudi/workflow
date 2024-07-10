package logic

type Code int

const (
	SUCCESS     Code = 0
	KeyExist    Code = 100
	KeyMiss     Code = 101
	SystemError Code = 500
)
