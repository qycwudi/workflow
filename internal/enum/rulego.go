package enum

// 记录运行状态 Space、Api
const (
	RecordStatusRunning = "running" // 运行中
	RecordStatusSuccess = "success" // 成功
	RecordStatusFail    = "fail"    // 失败
	RecordStatusCancel  = "cancel"  // 已取消
)

const (
	TraceStatusPending = "pending" // 初始化
	TraceStatusRunning = "running" // 运行中
	TraceStatusFinish  = "finish"  // 完成
)

const (
	CanvasMsg = "CANVAS_MSG" // 画布消息
	ApiMsg    = "API_MSG"    // API消息
)
