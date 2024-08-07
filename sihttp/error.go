package sihttp

import (
	"bytes"
	"fmt"
	"net/http"
)

type Error struct {
	Response *http.Response
	Body     []byte
}

func (e Error) Error() string {
	msg := bytes.Buffer{}
	if e.Response == nil {
		msg.WriteString("status: unknown")
	} else {
		msg.WriteString(fmt.Sprintf("status: %s", e.Response.Status))
	}

	if e.Body != nil {
		if msg.Len() > 0 {
			msg.WriteString(", ")
		}
		msg.WriteString(fmt.Sprintf("body: %s", e.Body))
	}
	return msg.String()
}

func (e Error) GetStatusCode(defaultStatusCode int) int {
	if e.Response != nil {
		return e.Response.StatusCode
	}
	return defaultStatusCode
}
func (e Error) GetStatus(defaultStatusCode int) string {
	if e.Response != nil {
		return e.Response.Status
	}
	return http.StatusText(defaultStatusCode)
}
