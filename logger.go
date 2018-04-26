package main

import "os"

// Logger is a generic interface for any logger. fileLogger below implements it
// to log to a file.
type Logger interface {
	Log(msg string) error
	Close() error
}

type fileLogger struct {
	logLocation string
	file        *os.File
}

func NewFileLogger(logLocation string) *fileLogger {
	return &fileLogger{logLocation: logLocation}
}

func (f *fileLogger) Open() error {
	if f.file == nil {
		file, err := os.OpenFile(f.logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		f.file = file
	}
	return nil
}

func (f *fileLogger) Log(msg string) error {
	if _, err := f.file.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

func (f *fileLogger) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}
