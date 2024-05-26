package sihttp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServer_NewServer(t *testing.T) {

	expected := []byte("hello world")
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	})

	m := http.NewServeMux()
	m.Handle("/", h)
	c := tls.Config{}
	s := NewServer(m, &c, ":63000", 30*time.Second, 30*time.Second)
	require.NotNil(t, s)

	// go func() {
	// 	err := s.Start()
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// }()
	// time.Sleep(2 * time.Second)
	// resBody, err := _get("http://127.0.0.1:63000/")
	// require.Nil(t, err)
	// require.EqualValues(t, expected, resBody)
	// fmt.Println(string(resBody))
	// require.Nil(t, s.Stop())
}

func TestServer_NewServerTls(t *testing.T) {

	expected := []byte("hello world")
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	})

	m := http.NewServeMux()
	m.Handle("/", h)
	c := tls.Config{}
	s := NewServerTls(m, &c, ":63000", 30*time.Second, 3000*time.Second, "", "")
	require.NotNil(t, s)

}

func TestServer_CreateTLSConfigMinTls(t *testing.T) {
	c := CreateTLSConfigMinTls(tls.VersionTLS13)
	require.NotNil(t, c)
	require.EqualValues(t, tls.VersionTLS13, c.MinVersion)
}

func _get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c := NewClient(_newStandardClient())
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)

}
