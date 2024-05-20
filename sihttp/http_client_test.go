package sihttp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/sicore"
	"github.com/wonksing/si/v2/siutils"
)

func Test_Client_setDefaultHeader(t *testing.T) {
	sc := http.Client{}

	defaultHeaders := map[string]string{
		"Content-type": "application/json",
	}
	c := NewClient(&sc, WithDefaultHeaders(defaultHeaders))

	req, err := http.NewRequest(http.MethodGet, "", nil)
	siutils.AssertNilFail(t, err)
	c.setDefaultHeader(req)

	assert.EqualValues(t, "application/json", req.Header.Get("Content-type"))
}

func Test_Client_appendRequestOption(t *testing.T) {
	sc := http.Client{}

	c := NewClient(&sc)

	c.appendRequestOption(WithBasicAuth("my-user", "my-password"))

	assert.EqualValues(t, 1, len(c.requestOpts))
}

func Test_Client_appendWriterOption(t *testing.T) {
	sc := http.Client{}

	c := NewClient(&sc)

	c.appendWriterOption(sicore.SetDefaultEncoder())

	assert.EqualValues(t, 1, len(c.writerOpts))
}

func Test_Client_appendReaderOption(t *testing.T) {
	sc := http.Client{}

	c := NewClient(&sc)

	c.appendReaderOption(sicore.SetJsonDecoder())

	assert.EqualValues(t, 1, len(c.readerOpts))
}

func Test_Client_request(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	// Standard http client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	tr := &http.Transport{
		MaxIdleConns:       300,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    tlsConfig,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}
	sc := NewStandardClient(time.Duration(30), tr)

	c := NewClient(sc,
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody, err := c.request(context.Background(), http.MethodGet, svr.URL, nil, nil, nil)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	var decodedResBody map[string]string
	err = c.requestDecode(context.Background(), http.MethodGet, svr.URL, nil, nil, nil, &decodedResBody)
	siutils.AssertNilFail(t, err)
	b, _ := json.Marshal(&decodedResBody)
	assert.EqualValues(t, expected, b)

}

func Test_setHeader(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	siutils.AssertNilFail(t, err)

	h := http.Header{}
	h.Set("Content-type", "application/json")
	setHeader(req, h)
	assert.EqualValues(t, "application/json", req.Header.Get("Content-type"))
}

func Test_setQueries(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	siutils.AssertNilFail(t, err)

	q := make(map[string]string)
	q["msg"] = "hello there"
	setQueries(req, q)
	assert.EqualValues(t, "hello there", req.URL.Query().Get("msg"))
}

func Test_Client_isRetryError(t *testing.T) {
	c := NewClient(nil)

	assert.EqualValues(t, false, c.isRetryError(nil))

	e := &Error{
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
		},
	}
	assert.EqualValues(t, true, c.isRetryError(e))

	assert.EqualValues(t, false, c.isRetryError(errors.New("just error")))
}
