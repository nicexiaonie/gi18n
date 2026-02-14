package gi18n

import "errors"

// 预定义错误
var (
	// ErrMessageNotFound 翻译消息不存在
	ErrMessageNotFound = errors.New("gi18n: message not found")
	// ErrInvalidFormat 无效的文件格式
	ErrInvalidFormat = errors.New("gi18n: invalid file format")
	// ErrEmptyID 空的消息 ID
	ErrEmptyID = errors.New("gi18n: empty message ID")
)

// MissPolicy 翻译缺失时的处理策略
type MissPolicy int

const (
	// MissReturnID 返回消息 ID 作为后备（默认行为）
	MissReturnID MissPolicy = iota
	// MissReturnEmpty 返回空字符串
	MissReturnEmpty
)

// Logger 日志接口，兼容 slog / zap / logrus 等日志库
//
// 使用示例（slog）:
//
//	gi18n.Init(&gi18n.Config{
//	    Logger: slog.Default(),
//	})
type Logger interface {
	Warn(msg string, args ...any)
}
