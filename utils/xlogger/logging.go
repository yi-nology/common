package xlogger

type Logger interface {
	Errorf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
}
