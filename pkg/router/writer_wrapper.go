package router

import "net/http"

type Writer struct {
	http.ResponseWriter
	wroteCode bool
	code      int
}

func WriterWrapper(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) {
	if !w.wroteCode {
		w.setCode(code)
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *Writer) Write(data []byte) (int, error) {
	if !w.wroteCode {
		w.setCode(http.StatusOK)
	}

	return w.ResponseWriter.Write(data) //nolint:wrapcheck
}

func (w *Writer) Code() int {
	return w.code
}

func (w *Writer) setCode(code int) {
	w.wroteCode = true
	w.code = code
}
