package codec

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/siutils"
)

func Test_EncodeJson(t *testing.T) {

	type _testStruct struct {
		Msg string `json:"msg"`
	}

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

	type _testStruct struct {
		Msg string `json:"msg"`
	}

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

func Test_DecodeJson(t *testing.T) {
	var v []byte = []byte(`{"id":1,"email_address":"asdf","name":"asdf","borrowed":false,"book_id":23}`)
	buf := bytes.NewBuffer(v)
	out := make(map[string]interface{})

	err := DecodeJson(&out, buf)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "asdf", out["email_address"].(string))
}

func Test_DecodeJsonCopied(t *testing.T) {
	var v []byte = []byte(`{"id":1,"email_address":"asdf","name":"asdf","borrowed":false,"book_id":23}`)
	buf := bytes.NewBuffer(v)
	out := make(map[string]interface{})

	copied, err := DecodeJsonCopied(&out, buf)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, "asdf", out["email_address"].(string))
	assert.EqualValues(t, v, copied.Bytes())
}
