package logger

type Loggerer interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Loggerer
}
