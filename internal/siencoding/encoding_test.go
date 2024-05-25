package siencoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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
