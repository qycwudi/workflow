package utils

import "github.com/zeromicro/go-zero/core/logx"

// RoleCustomLog 自定义日志输出
type RoleCustomLog struct {
}

func (l *RoleCustomLog) Printf(format string, v ...interface{}) {
	logx.Infof(format, v...)
}
