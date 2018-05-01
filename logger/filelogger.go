package logger

import "os"

type FileLogger struct {
	logLocation string
	file        *os.File
}

func NewFileLogger(logLocation string) *FileLogger {
	return &FileLogger{logLocation: logLocation}
}

func (f *FileLogger) Open() error {
	if f.file == nil {
		file, err := os.OpenFile(f.logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		f.file = file
	}
	return nil
}

func (f *FileLogger) Log(msg string) error {
	if _, err := f.file.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

func (f *FileLogger) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}
