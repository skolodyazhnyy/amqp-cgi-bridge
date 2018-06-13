package log

import (
	"fmt"
	"io"
)

func NewStdLogger(logger *Logger) *StdLogger {
	return &StdLogger{
		logger: logger,
	}
}

type StdLogger struct {
	logger *Logger
}

func (l *StdLogger) SetOutput(w io.Writer) {
}

func (l *StdLogger) Output(calldepth int, s string) error {
	return nil
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *StdLogger) Print(v ...interface{}) {
	l.logger.Debug(fmt.Sprint(v...), nil)
}

func (l *StdLogger) Println(v ...interface{}) {
	l.logger.Debug(fmt.Sprint(v...), nil)
}

func (l *StdLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *StdLogger) Fatalln(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *StdLogger) Panic(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *StdLogger) Panicf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *StdLogger) Panicln(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *StdLogger) Flags() int {
	return 0
}

func (l *StdLogger) SetFlags(flag int) {
}

func (l *StdLogger) Prefix() string {
	return ""
}

func (l *StdLogger) SetPrefix(prefix string) {
}
