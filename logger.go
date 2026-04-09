package anrok

// Logger 日志记录器接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
	Key   string
	Value any
}

// String 创建字符串字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Number 创建数字字段
func Number(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// ErrorField 创建错误字段
func ErrorField(err error) Field {
	return Field{Key: "error", Value: err}
}

// NopLogger 空日志记录器（不记录任何日志）
type NopLogger struct{}

func (n *NopLogger) Debug(msg string, fields ...Field) {}
func (n *NopLogger) Info(msg string, fields ...Field)  {}
func (n *NopLogger) Error(msg string, fields ...Field) {}
