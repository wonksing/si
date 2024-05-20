package sicore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/siutils"
)

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

func Test_DefaultDecoder(t *testing.T) {
	d := NewDefaultDecoder(nil)
	siutils.AssertNotNilFail(t, d)

	var v []byte = []byte(`{"id":1,"email_address":"asdf","name":"asdf","borrowed":false,"book_id":23}`)
	buf := bytes.NewBuffer(v)

	d.Reset(buf)

	var err error
	var out []byte

	err = d.Decode(&out)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, v, out)

	buf = bytes.NewBuffer(v)
	d.Reset(buf)
	var outStr string
	err = d.Decode(&outStr)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, string(v), outStr)
}
