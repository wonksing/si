package siencoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/siutils"
)

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
