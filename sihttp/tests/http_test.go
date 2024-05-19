package sihttp_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/sicore"
	"github.com/wonksing/si/v2/sihttp"
	"github.com/wonksing/si/v2/siutils"
	"github.com/wonksing/si/v2/tests/testmodels"
)

func Test_Online_Client_Do(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, standardClient)

	hc := sihttp.NewClient(standardClient)

	request, err := http.NewRequest(http.MethodGet, remoteAddr+"/test/hello", nil)
	siutils.AssertNilFail(t, err)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(request)
	siutils.AssertNilFail(t, err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(b))
}

func Test_Client_Post(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	data := "hey"
	url := remoteAddr + "/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.Post(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))
}

func Test_Client_Post_inputReader(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	data := "hey"
	url := remoteAddr + "/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.Post(url, nil, strings.NewReader(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))
}

func Test_Client_Post_fileData(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	url := remoteAddr + "/test/echo"

	f, err := os.OpenFile("./data/testfile.txt", os.O_RDONLY, 0777)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	header := make(http.Header)
	header["Content-Type"] = []string{"multipart/form-data"}

	res, err := client.Post(url, header, f)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestCheckRequestState(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithReadWriter(buf)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	req.Header.Set("custom_header", "wonk")

	sendData := fmt.Sprintf("%s-%d", data, 0)
	rw.WriteFlush([]byte(sendData))
	resp, err := standardClient.Do(req)
	siutils.AssertNilFail(t, err)

	respBody, err := io.ReadAll(resp.Body)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))
	resp.Body.Close()

	req2, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	for k := range req.Header {
		delete(req.Header, k)
	}

	assert.EqualValues(t, req2, req)
}
func TestReuseRequest(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithReadWriter(buf)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	req.Header.Set("custom_header", "wonk")

	for i := 0; i < 10; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)
		rw.WriteFlush([]byte(sendData))
		resp, err := standardClient.Do(req)
		siutils.AssertNilFail(t, err)

		respBody, err := io.ReadAll(resp.Body)
		siutils.AssertNilFail(t, err)
		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))

		resp.Body.Close()
	}

}

func TestReuseRequestInGoroutinePanic(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	t.Skip("skipping because this code panics")
	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithReadWriter(buf)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	var wg sync.WaitGroup
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for i := 0; i < 10; i++ {
				sendData := fmt.Sprintf("%s-%d", data, i)

				req.Header.Set("custom_header", sendData)

				rw.WriteFlush([]byte(sendData))
				resp, err := standardClient.Do(req)
				siutils.AssertNilFail(t, err)

				respBody, err := io.ReadAll(resp.Body)
				siutils.AssertNilFail(t, err)
				assert.EqualValues(t, sendData, string(respBody))
				fmt.Println(string(respBody))

				resp.Body.Close()
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()

}

func TestReuseRequestInGoroutine(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	var wg sync.WaitGroup
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, routineNumber int) {
			buf := bytes.NewBuffer(make([]byte, 0, 1024))
			rw := sicore.GetReadWriterWithReadWriter(buf)

			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", nil)
			siutils.AssertNilFail(t, err)

			req.Body = ioutil.NopCloser(rw)

			for i := 0; i < 10; i++ {
				sendData := fmt.Sprintf("%s-%d-%d", data, routineNumber, i)

				req.Header.Set("custom_header", sendData)

				rw.WriteFlush([]byte(sendData))
				resp, err := standardClient.Do(req)
				siutils.AssertNilFail(t, err)

				respBody, err := io.ReadAll(resp.Body)
				siutils.AssertNilFail(t, err)
				assert.EqualValues(t, sendData, string(respBody))
				fmt.Println(string(respBody))

				resp.Body.Close()
			}

			sicore.PutReadWriter(rw)
			wg.Done()
		}(&wg, j)
	}
	wg.Wait()

}

func TestHttpClientRequestPostTls(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	data := "hey"
	urls := []string{"http://127.0.0.1:8080/test/echo", "https://127.0.0.1:8081/test/echo"}
	for i := 0; i < 2; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)

		respBody, err := client.Post(urls[i], nil, []byte(sendData))
		siutils.AssertNilFail(t, err)

		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))
	}
}

func TestHttpClientRequestGet(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient, sihttp.WithWriterOpt(sicore.SetJsonEncoder()))

	url := "http://127.0.0.1:8080/test/hello"

	queries := make(map[string]string)
	queries["name"] = "wonk"
	queries["kor"] = "길동"

	respBody, err := client.Get(url, nil, queries)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(respBody))
	// fmt.Println(string(respBody))

}
func TestHttpClientRequestPost(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.Post(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPut(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.Put(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPostJsonDecoded(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()))

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestHttpClientRequestPostReaderFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient)

	url := "http://127.0.0.1:8080/test/file/upload"

	f, err := os.OpenFile("./data/testfile.txt", os.O_RDONLY, 0777)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	contents, err := io.ReadAll(f)
	siutils.AssertNilFail(t, err)

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	mw := multipart.NewWriter(buf)

	part, err := mw.CreateFormFile("file_to_upload", f.Name())
	siutils.AssertNilFail(t, err)
	part.Write(contents)

	mw.WriteField("nam", "wonk")

	header := make(http.Header)
	header["Content-Type"] = []string{mw.FormDataContentType()}

	err = mw.Close()
	siutils.AssertNilFail(t, err)

	// res, err := client.RequestPostFile(url, header, buf)
	res, err := client.Post(url, header, buf)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestHttpClientRequestPostFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmacKey", []byte("1234"))),
	)

	url := "http://127.0.0.1:8080/test/file/upload"

	res, err := client.PostFile(url, nil, nil, "file_to_upload", "./data/testfile.txt")
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestHttpClientRequestGetWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmac-hash", []byte("1234"))),
	)

	url := "http://127.0.0.1:8080/test/hello"

	respBody, err := client.Get(url, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(respBody))

}
func TestHttpClientRequestPostWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmac-hash", []byte("1234"))),
	)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.Post(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPostJsonDecodedWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithRequestHeaderHmac256("hmacKey", []byte("1234")),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	// client.SetRequestOptions(sihttp.WithHeaderHmac256("hmacKey", []byte("1234")))
	// client.SetWriterOptions(sicore.SetJsonEncoder())
	// client.SetReaderOptions(sicore.SetJsonDecoder())

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestHttpClientRequestPostJsonDecodedWithBearerToken(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithRequestOpt(sihttp.WithBearerToken("asdf")),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.PostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestWithBaseUrl(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(standardClient,
		sihttp.WithBaseUrl("http://127.0.0.1:8080"),
	)

	url := "/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	b, _ := json.Marshal(&student)
	res, err := client.Post(url, nil, b)
	siutils.AssertNilFail(t, err)

	expected := `{"id":1,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false}`
	assert.EqualValues(t, expected, string(res))

}

func Test_Client_NewClient(t *testing.T) {
	c := sihttp.NewClient(newStandardClient())
	siutils.AssertNotNilFail(t, c)
}

func Test_Client_Do(t *testing.T) {

	expected := []byte("hello there")

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient())
	siutils.AssertNotNilFail(t, c)

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	siutils.AssertNilFail(t, err)
	resp, err := c.Do(req)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_DoRead(t *testing.T) {

	expected := []byte("hello there")

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient())
	siutils.AssertNotNilFail(t, c)

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	siutils.AssertNilFail(t, err)
	resBody, err := c.DoRead(req)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_DoDecode(t *testing.T) {

	expected := []byte("hello there")

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient())
	siutils.AssertNotNilFail(t, c)

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	siutils.AssertNilFail(t, err)

	var resBody []byte
	err = c.DoDecode(req, &resBody)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_DoDecode_Struct(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient(),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	siutils.AssertNilFail(t, err)

	resBody := testStruct{}
	err = c.DoDecode(req, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Request(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient(),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody, err := c.Request(http.MethodGet, svr.URL, nil, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_RequestContext(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient(),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody, err := c.RequestContext(context.Background(), http.MethodGet, svr.URL, nil, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_RequestDecode(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient(),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody := testStruct{}
	err := c.RequestDecode(http.MethodGet, svr.URL, nil, nil, nil, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_RequestDecodeContext(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := sihttp.NewClient(newStandardClient(),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody := testStruct{}
	err := c.RequestDecodeContext(context.Background(), http.MethodGet, svr.URL, nil, nil, nil, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func newStandardClient() *http.Client {
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

	return sihttp.NewStandardClient(time.Duration(30), tr)
}

type testStruct struct {
	Msg string `json:"msg"`
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
