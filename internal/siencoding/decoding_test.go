package siencoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DefaultDecoder(t *testing.T) {
	d := NewDefaultDecoder(nil)
	require.NotNil(t, d)

	var v []byte = []byte(`{"id":1,"email_address":"asdf","name":"asdf","borrowed":false,"book_id":23}`)
	buf := bytes.NewBuffer(v)

	d.Reset(buf)

	var err error
	var out []byte

	err = d.Decode(&out)
	require.Nil(t, err)
	assert.EqualValues(t, v, out)

	buf = bytes.NewBuffer(v)
	d.Reset(buf)
	var outStr string
	err = d.Decode(&outStr)
	require.Nil(t, err)
	assert.EqualValues(t, string(v), outStr)
}
