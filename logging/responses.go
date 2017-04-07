package logging

import "net/http"

// ResponseLogger records http responses
type ResponseLogger struct {
	Writer http.ResponseWriter
	status int
	size   int
}

// Header returns the headers
func (l *ResponseLogger) Header() http.Header {
	return l.Writer.Header()
}

// Write write the body
func (l *ResponseLogger) Write(b []byte) (int, error) {
	size, err := l.Writer.Write(b)
	l.size += size
	return size, err
}

// WriteHeader write the status
func (l *ResponseLogger) WriteHeader(s int) {
	l.Writer.WriteHeader(s)
	l.status = s
}

// Status returns the status
func (l *ResponseLogger) Status() int {
	if l.status == 0 {
		return http.StatusOK
	}
	return l.status
}

// Size returns the response body size
func (l *ResponseLogger) Size() int {
	return l.size
}

func (l *ResponseLogger) Flush() {
	l.Writer.(http.Flusher).Flush()
}
