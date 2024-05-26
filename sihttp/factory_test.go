package sihttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_DefaultInsecureStandardClient(t *testing.T) {
	c := DefaultInsecureStandardClient()
	require.NotNil(t, c)
}

func Test_DefaultStandardClient(t *testing.T) {
	config := tls.Config{}
	c := DefaultStandardClient(&config)
	require.NotNil(t, c)
}

func Test_NewStandardClient(t *testing.T) {
	config := &tls.Config{}

	dialer := &net.Dialer{Timeout: 5 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       50,
		IdleConnTimeout:    time.Duration(60) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    config,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	c := NewStandardClient(30*time.Second, tr)
	require.NotNil(t, c)
}
