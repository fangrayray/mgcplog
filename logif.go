package mgcplog

// LogInterface : Log interface
type LogInterface interface {
	Info(le LogEntity)
	Warn(le LogEntity)
	Error(le LogEntity)
	Fetal(le LogEntity)
	Panic(le LogEntity) error
}
