package internal

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EncodeJson(t *testing.T) {
	dst := bytes.Buffer{}
	src := _testStruct{
		Msg: "hello world",
	}
	err := EncodeJson(&dst, &src)
	require.Nil(t, err)
	b, _ := json.Marshal(&src)
	b = append(b, '\n')
	require.EqualValues(t, b, dst.Bytes())
}

func Test_EncodeJsonCopied(t *testing.T) {
	dst := bytes.Buffer{}
	src := _testStruct{
		Msg: "hello world",
	}
	copied, err := EncodeJsonCopied(&dst, &src)
	require.Nil(t, err)
	b, _ := json.Marshal(&src)
	b = append(b, '\n')
	require.EqualValues(t, b, dst.Bytes())
	require.EqualValues(t, b, copied.Bytes())
}

func Test_HmacSha256HexEncoded(t *testing.T) {
	res, err := HmacSha256HexEncoded("asdf", []byte("hello world"))
	require.Nil(t, err)
	require.EqualValues(t, "2c78fedf60d1f955bf0c9e14ed6b332a6efb6e5668fc6aa067257558cdbb7d6d", res)
}

func Test_HmacSha256HexEncodedWithReader(t *testing.T) {
	r := bytes.NewBufferString("hello world")
	res, err := HmacSha256HexEncodedWithReader("asdf", r)
	require.Nil(t, err)
	require.EqualValues(t, "2c78fedf60d1f955bf0c9e14ed6b332a6efb6e5668fc6aa067257558cdbb7d6d", res)
}

func TestDefaultEncoder_Reset(t *testing.T) {
	b := &bytes.Buffer{}
	e := DefaultEncoder{}

	e.Reset(b)
	require.Equal(t, b, e.w)
}

func TestDefaultEncoder_Encode(t *testing.T) {

	t.Run("succeed-bytes", func(t *testing.T) {
		b := &bytes.Buffer{}
		e := DefaultEncoder{
			w: b,
		}

		data := []byte("hello")
		err := e.Encode(data)
		require.Nil(t, err)

		require.EqualValues(t, data, b.Bytes())
	})

	t.Run("succeed-bytes_ptr", func(t *testing.T) {
		b := &bytes.Buffer{}
		e := DefaultEncoder{
			w: b,
		}

		data := []byte("hello")
		err := e.Encode(&data)
		require.Nil(t, err)

		require.EqualValues(t, data, b.Bytes())
	})
	t.Run("succeed-string", func(t *testing.T) {
		b := &bytes.Buffer{}
		e := DefaultEncoder{
			w: b,
		}

		data := "hello"
		err := e.Encode(data)
		require.Nil(t, err)

		require.EqualValues(t, data, b.String())
	})
	t.Run("succeed-string_ptr", func(t *testing.T) {
		b := &bytes.Buffer{}
		e := DefaultEncoder{
			w: b,
		}

		data := "hello"
		err := e.Encode(&data)
		require.Nil(t, err)

		require.EqualValues(t, data, b.String())
	})
	t.Run("fail", func(t *testing.T) {
		b := &bytes.Buffer{}
		e := DefaultEncoder{
			w: b,
		}

		data := 1535
		err := e.Encode(&data)
		require.NotNil(t, err)
	})
}

type _testStruct struct {
	Msg string `json:"msg"`
}
