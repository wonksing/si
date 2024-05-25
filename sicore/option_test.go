package sicore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriterOptionFunc_apply(t *testing.T) {
	var f WriterOptionFunc = func(w *Writer) {}

	buf := bytes.Buffer{}
	w := newWriter(&buf)
	f.apply(w)
}

func Test_SetJsonEncoder(t *testing.T) {
	w := Writer{}
	o := SetJsonEncoder()
	o.apply(&w)
	require.NotNil(t, w.enc)
}

func Test_SetDefaultEncoder(t *testing.T) {
	w := Writer{}
	o := SetDefaultEncoder()
	o.apply(&w)
	require.NotNil(t, w.enc)
}

func Test_SetEofChecker(t *testing.T) {
	c := DefaultEofChecker
	o := SetEofChecker(&c)

	r := Reader{}
	o.apply(&r)
	require.Equal(t, &c, r.chk)
}

func Test_SetDefaultEOFChecker(t *testing.T) {
	o := SetDefaultEOFChecker()

	r := Reader{}
	o.apply(&r)

	c := DefaultEofChecker
	require.Equal(t, &c, r.chk)
}

func Test_SetJsonDecoder(t *testing.T) {
	o := SetJsonDecoder()

	r := Reader{}

	o.apply(&r)
	require.NotNil(t, r.dec)
}

func Test_WithTagKey(t *testing.T) {
	o := WithTagKey("si")

	r := RowScanner{}

	o.apply(&r)

	require.EqualValues(t, "si", r.tagKey)
}
