package event

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	contentType = []string{"text/event-stream"}
	noCache     = []string{"no-cache"}
)

type stringWriter interface {
	io.Writer
	writeString(string) (int, error)
}

type stringWrapper struct {
	io.Writer
}

func (w stringWrapper) writeString(str string) (int, error) {
	return w.Writer.Write([]byte(str))
}

func checkWriter(writer io.Writer) stringWriter {
	if w, ok := writer.(stringWriter); ok {
		return w
	} else {
		return stringWrapper{writer}
	}
}

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

type Event struct {
	Data interface{} `json:"data"`
}

func (r Event) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return encode(w, r)
}

func (r Event) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = contentType

	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = noCache
	}
}

func encode(writer io.Writer, event Event) error {
	w := checkWriter(writer)
	return writeData(w, event.Data)
}

func writeData(w stringWriter, data interface{}) error {
	_, _ = dataReplacer.WriteString(w, fmt.Sprint(data))
	if strings.HasPrefix(data.(string), "data") {
		_, _ = w.writeString("\n\n")
	}
	return nil
}
