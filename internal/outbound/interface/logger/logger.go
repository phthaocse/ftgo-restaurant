package logger

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}
