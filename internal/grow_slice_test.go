package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_grow(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		b := make([]byte, 0, 10)
		b = append(b, []byte("asdfㅁ")...)
		assert.Equal(t, 7, len(b))
		assert.Equal(t, 10, cap(b))

		l, err := grow(&b, 100)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, 7, l)
		assert.Equal(t, 107, len(b))
		assert.Equal(t, 120, cap(b))
	})

	t.Run("succeed-available", func(t *testing.T) {
		b := make([]byte, 0, 100)
		b = append(b, []byte("0123456789")...)

		l, err := grow(&b, 20)
		require.Nil(t, err)
		require.EqualValues(t, 10, l)
		require.EqualValues(t, 30, len(b))
		require.EqualValues(t, 100, cap(b))
	})

	t.Run("fail-too_large", func(t *testing.T) {
		b := make([]byte, 0, 100)
		b = append(b, []byte("0123456789")...)

		_, err := grow(&b, maxInt)
		require.NotNil(t, err)
	})
}

func Test_growCap(t *testing.T) {

	t.Run("succeed", func(t *testing.T) {
		b := make([]byte, 0, 10)
		b = append(b, []byte("asdfㅁ")...)
		assert.Equal(t, 7, len(b))
		assert.Equal(t, 10, cap(b))

		err := growCap(&b, 100)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, 7, len(b))
		assert.Equal(t, 110, cap(b))
	})

	t.Run("succeed", func(t *testing.T) {
		b := make([]byte, 0, 100)
		b = append(b, []byte("0123456789")...)

		err := growCap(&b, 100)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, 10, len(b))
		assert.Equal(t, 200, cap(b))
	})

	t.Run("succeed-available", func(t *testing.T) {
		b := make([]byte, 0, 100)
		b = append(b, []byte("0123456789")...)

		err := growCap(&b, 50)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, 10, len(b))
		assert.Equal(t, 100, cap(b))
	})
	t.Run("fail-too_large", func(t *testing.T) {
		b := make([]byte, 0, 100)
		b = append(b, []byte("0123456789")...)

		err := growCap(&b, maxInt)
		require.NotNil(t, err)
	})
}
