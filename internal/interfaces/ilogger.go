package interfaces

type ILogger interface {
	Info(msg string, v ...interface{})
	Error(err error)
	ErrorWithMsg(msg string, err error)
}
