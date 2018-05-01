package logger

// Logger is a generic interface for any logger. FileLogger below implements it
// to log to a file.
type Logger interface {
	Log(msg string) error
	Close() error
}
