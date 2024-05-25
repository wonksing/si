package sio

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

func Test_WithSqlColumnType(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"
		o := WithSqlColumnType(name, SqlColTypeBool)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullBoolTypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("byte", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeByte)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullByteTypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("bytes", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeBytes)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlBytesTypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("string", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeString)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullStringTypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("int", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeInt)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullIntTypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("int8", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeInt8)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullInt8TypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("int16", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeInt16)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullInt16TypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("int32", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeInt32)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullInt32TypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("int64", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeInt64)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullInt64TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uint", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUint)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullUintTypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uint8", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUint8)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullUint8TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uint16", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUint16)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullUint16TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uint32", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUint32)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullUint32TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uint64", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUint64)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullUint64TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("float32", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeFloat32)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullFloat32TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("float64", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeFloat64)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullFloat64TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("time", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeTime)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, sqlNullTimeTypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("ints", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeints)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, intsTypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("ints8", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeints8)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, ints8TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("ints16", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeints16)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, ints16TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("ints32", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeints32)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, ints32TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("ints64", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeints64)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, ints64TypeValue)
		} else {
			t.FailNow()
		}
	})

	t.Run("uints", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUints)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, uintsTypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uints8", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUints8)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, uints8TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uints16", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUints16)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, uints16TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uints32", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUints32)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, uints32TypeValue)
		} else {
			t.FailNow()
		}
	})
	t.Run("uints64", func(t *testing.T) {
		r := RowScanner{
			sqlCol: make(map[string]any),
		}
		name := "some_key_name"

		o := WithSqlColumnType(name, SqlColTypeUints64)
		o.apply(&r)
		if a, ok := r.sqlCol[name]; ok {
			require.EqualValues(t, a, uints64TypeValue)
		} else {
			t.FailNow()
		}
	})
}
