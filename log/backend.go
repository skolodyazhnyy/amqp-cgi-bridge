package log

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type backend interface {
	Record(R)
}

func NewX(fmt string, w io.Writer, f func(R) string) *Logger {
	switch fmt {
	case "json":
		return NewJSON(w)
	default:
		return NewText(w, f)
	}
}

func NewText(w io.Writer, f func(R) string) *Logger {
	return New(NewTextBackend(w, f))
}

func NewJSON(w io.Writer) *Logger {
	return New(NewJSONBackend(w))
}

type JSONBackend struct {
	writer io.Writer
}

func NewJSONBackend(w io.Writer) *JSONBackend {
	return &JSONBackend{writer: w}
}

func (b *JSONBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	rec["timestamp"] = time.Now().UTC()

	json.NewEncoder(b.writer).Encode(rec)
}

type TextBackend struct {
	writer    io.Writer
	formatter func(R) string
}

func DefaultTextFormat(rec R) string {
	channel := rec["channel"]
	if channel == nil {
		channel = "main"
	}
	delete(rec, "channel")

	severity := rec["severity"]
	if severity == nil {
		severity = "?"
	}
	delete(rec, "severity")

	message := rec["message"]
	delete(rec, "message")

	data, err := json.Marshal(rec)
	if err != nil {
		data = []byte(fmt.Sprintf("[!marshalling error: %v!]", err))
	}

	return fmt.Sprintf("%v %v.%v %v %s", time.Now().Format(time.RFC3339Nano), channel, severity, message, data)
}

func NewTextBackend(w io.Writer, f func(R) string) *TextBackend {
	if f == nil {
		f = DefaultTextFormat
	}

	return &TextBackend{writer: w, formatter: f}
}

func (b *TextBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	fmt.Fprintln(b.writer, b.formatter(rec))
}
