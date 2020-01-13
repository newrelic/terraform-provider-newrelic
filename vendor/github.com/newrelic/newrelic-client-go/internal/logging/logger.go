package logging

// Logger interface implements a simple logger.
type Logger interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Warn(string, ...interface{})
}
