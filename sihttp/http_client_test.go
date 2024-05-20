package sihttp

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/sicore"
	"github.com/wonksing/si/v2/siutils"
)

func Test_Client_NewClient(t *testing.T) {
	sc := _newStandardClient()
	c := NewClient(sc)
	siutils.AssertNotNilFail(t, c)

	c = NewClient(sc, nil)
	siutils.AssertNotNilFail(t, c)
}

func Test_Client_Do(t *testing.T) {

	expected := []byte("hello there")

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient())
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

	c := NewClient(_newStandardClient())
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

	c := NewClient(_newStandardClient())
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

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	siutils.AssertNilFail(t, err)

	resBody := _testStruct{}
	err = c.DoDecode(req, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := _jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Request(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	headers := http.Header{}
	queries := make(map[string]string)
	var reqBody []byte

	resBody, err := c.Request(http.MethodGet, svr.URL, headers, queries, reqBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resBody, err = c.RequestContext(context.Background(), http.MethodGet, svr.URL, headers, queries, reqBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)
}

func Test_Client_RequestDecode(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody := _testStruct{}
	err := c.RequestDecode(http.MethodGet, svr.URL, nil, nil, nil, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := _jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_RequestDecodeContext(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody := _testStruct{}
	err := c.RequestDecodeContext(context.Background(), http.MethodGet, svr.URL, nil, nil, nil, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := _jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Get(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody, err := c.Get(svr.URL, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_GetContext(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody, err := c.GetContext(context.Background(), svr.URL, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, expected, resBody)
}

func Test_Client_GetDecode(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	resBody := _testStruct{}
	err := c.GetDecode(svr.URL, nil, nil, &resBody)
	siutils.AssertNilFail(t, err)

	b, err := _jsonMarshal(&resBody)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)

	resBody2 := _testStruct{}
	err = c.GetDecodeContext(context.Background(), svr.URL, nil, nil, &resBody2)
	siutils.AssertNilFail(t, err)

	b2, err := _jsonMarshal(&resBody2)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b2)
}

func Test_Client_Post(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		w.Write(b)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	var err error
	var resBody []byte
	var resStruct _testStruct

	resBody, err = c.Post(svr.URL, nil, expected)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resBody, err = c.PostContext(context.Background(), svr.URL, nil, expected)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resStruct = _testStruct{}
	err = c.PostDecode(svr.URL, nil, expected, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err := _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)

	resStruct = _testStruct{}
	err = c.PostDecodeContext(context.Background(), svr.URL, nil, expected, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err = _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Put(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	var err error
	var resBody []byte
	var resStruct _testStruct

	resBody, err = c.Put(svr.URL, nil, nil)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resBody, err = c.PutContext(context.Background(), svr.URL, nil, nil)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resStruct = _testStruct{}
	err = c.PutDecode(svr.URL, nil, nil, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err := _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)

	resStruct = _testStruct{}
	err = c.PutDecodeContext(context.Background(), svr.URL, nil, nil, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err = _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Patch(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		w.Write(b)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	var err error
	var resBody []byte
	var resStruct _testStruct

	resBody, err = c.Patch(svr.URL, nil, expected)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resBody, err = c.PatchContext(context.Background(), svr.URL, nil, expected)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resStruct = _testStruct{}
	err = c.PatchDecode(svr.URL, nil, expected, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err := _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)

	resStruct = _testStruct{}
	err = c.PatchDecodeContext(context.Background(), svr.URL, nil, expected, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err = _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_Delete(t *testing.T) {

	expected := []byte(`{"msg":"hello there"}`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query().Get("msg")
		s := _testStruct{
			Msg: v,
		}
		res, _ := _jsonMarshal(&s)
		w.Write(res)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	var err error
	var resBody []byte
	var resStruct _testStruct
	queries := make(map[string]string)
	queries["msg"] = "hello there"

	resBody, err = c.Delete(svr.URL, nil, queries)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resBody, err = c.DeleteContext(context.Background(), svr.URL, nil, queries)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, resBody)

	resStruct = _testStruct{}
	err = c.DeleteDecode(svr.URL, nil, queries, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err := _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)

	resStruct = _testStruct{}
	err = c.DeleteDecodeContext(context.Background(), svr.URL, nil, queries, &resStruct)
	siutils.AssertNilFail(t, err)
	b, err = _jsonMarshal(&resStruct)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, b)
}

func Test_Client_PostFile(t *testing.T) {

	expected := []byte(`success`)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := "./tests/data/upload/"
		var err error

		// multipart/form-data로 파싱된 요청본문을 최대 1메가까지 메모리에 저장하도록 한다.
		// r.ParseMultipartForm(1 << 20)
		r.ParseMultipartForm(1 * 1024)

		// FormFile returns the first file for the provided form key.
		// FormFile calls ParseMultipartForm and ParseForm if necessary.
		// 첫번째 파일 데이터와 헤더를 반환한다. ParseMultipartForm과 ParseForm을 호출할 수 있다는데 언제인지는 모르겠다.
		file, header, err := r.FormFile("file_to_upload")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Uploaded File: %+v, File Size: %+v, MIME Header: %+v\n",
			header.Filename, header.Size, header.Header)

		// filePath 디렉토리가 없으면 만들기
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// 경로와 파일명 붙이기
		filePathName := filepath.Join(filePath, header.Filename)

		// 파일 만들기
		f, err := os.Create(filePathName)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer f.Close()

		// 멀티파트 파일 받아서 읽기 위함
		reader := bufio.NewReader(file)

		// 어디까지 읽었는지 보기 위함, 결국엔 사이즈랑 같아야 함
		var offset int64 = 0

		// reader로부터 4096 바이트씩 읽을 것임
		rb := make([]byte, 4096)
		for {
			size, err := reader.Read(rb) // rb에 집어넣기
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// n, err := f.WriteAt(rb[:size], offset)
			n, err := f.Write(rb[:size])
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			offset += int64(n)
		}
		log.Printf("file size: %v, %v", header.Size, offset)
		w.Write(expected)
	}))
	defer svr.Close()

	c := NewClient(_newStandardClient(),
		WithWriterOpt(sicore.SetJsonEncoder()),
		WithReaderOpt(sicore.SetJsonDecoder()),
	)
	siutils.AssertNotNilFail(t, c)

	var err error

	res, err := c.PostFile(svr.URL, nil, nil, "file_to_upload", "./tests/data/testfile.txt")
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, res)

	res, err = c.PostFileContext(context.Background(), svr.URL, nil, nil, "file_to_upload", "./tests/data/testfile.txt")
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, expected, res)
}

func _newStandardClient() *http.Client {
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

	return NewStandardClient(time.Duration(30), tr)
}

type _testStruct struct {
	Msg string `json:"msg"`
}

func _jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

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
