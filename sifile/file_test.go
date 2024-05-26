package sifile

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sio"
)

func TestNewFile(t *testing.T) {
	f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0777)
	require.Nil(t, err)
	defer f.Close()

	byt, err := f.ReadAll()
	require.Nil(t, err)
	fmt.Println(string(byt) + "1")

	_, err = f.WriteFlush([]byte("hey\n"))
	require.Nil(t, err)

	byt, err = f.ReadAllAt(0)
	require.Nil(t, err)
	fmt.Println(string(byt) + "2")

	_, err = f.WriteFlush([]byte("hey2\n"))
	require.Nil(t, err)

	byt, err = f.ReadAllAt(0)
	require.Nil(t, err)
	fmt.Println(string(byt) + "3")

	_, err = f.WriteFlush([]byte("hey3\n"))
	require.Nil(t, err)

	byt, err = f.ReadAllAt(0)
	require.Nil(t, err)
	fmt.Println(string(byt) + "4")

	err = f.Chdir()
	require.NotNil(t, err)

	err = f.Chmod(0777)
	require.Nil(t, err)

	require.NotEmpty(t, f.Name())

	t.Run("succeed", func(t *testing.T) {

		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0777)
		require.Nil(t, err)

		_, err = f.WriteFlush([]byte("hey2\n"))
		if err != nil {
			f.Close()
			require.Nil(t, err)
		}
		f.Close()

		f, err = OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		require.Nil(t, err)
		defer f.Close()

		p := make([]byte, 100)
		n, err := f.Read(p)
		require.Nil(t, err)
		require.EqualValues(t, 5, n)
	})

	t.Run("succeed", func(t *testing.T) {

		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		require.Nil(t, err)

		_, err = f.WriteAt([]byte("hey2\n"), 0)
		if err != nil {
			f.Close()
			require.Nil(t, err)
		}
		f.Close()

		f, err = OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		require.Nil(t, err)
		defer f.Close()

		p := make([]byte, 100)
		n, _ := f.ReadAt(p, 2)
		require.EqualValues(t, 3, n)
	})

	t.Run("succeed", func(t *testing.T) {

		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		require.Nil(t, err)

		_, err = f.WriteString("hey2\n")
		if err != nil {
			f.Close()
			require.Nil(t, err)
		}
		require.Nil(t, f.Flush())
		f.Close()

		f, err = OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		require.Nil(t, err)
		defer f.Close()

		p := make([]byte, 100)
		n, _ := f.ReadAt(p, 2)
		require.EqualValues(t, 3, n)
	})
}

func Test_OpenFile(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		_, err := OpenFile("", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0777)
		require.NotNil(t, err)
	})
}

func Test_Create(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		_, err := Create("./tests/data/test.txt",
			WithReaderOpt(sio.SetJsonDecoder()),
			WithWriterOpt(sio.SetJsonEncoder()))
		require.Nil(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		_, err := Create("")
		require.NotNil(t, err)
	})
}

func TestFile_Read(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		f, err := Open("./tests")
		require.Nil(t, err)
		defer f.Close()

		res, err := f.ReadDir(2)
		require.Nil(t, err)
		require.EqualValues(t, 2, len(res))

		res2, err := f.Readdir(1)
		require.Nil(t, err)
		require.EqualValues(t, 1, len(res2))

	})
	t.Run("succeed", func(t *testing.T) {
		f, err := Open("./tests")
		require.Nil(t, err)
		defer f.Close()

		res3, err := f.Readdirnames(1)
		require.Nil(t, err)
		require.EqualValues(t, 1, len(res3))
	})

	t.Run("succeed-2", func(t *testing.T) {

		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		require.Nil(t, err)
		defer f.Close()

		data := bytes.NewBufferString("hello world")
		n, err := f.ReadFrom(data)
		require.Nil(t, err)
		require.EqualValues(t, 11, n)
	})

	t.Run("succeed-3", func(t *testing.T) {
		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
		require.Nil(t, err)
		defer f.Close()

		data := []byte("hello world")
		err = f.Encode(data)
		require.Nil(t, err)
		f.Flush()

		data = []byte("hello world")
		err = f.EncodeFlush(data)
		require.Nil(t, err)

	})

	t.Run("succeed-4", func(t *testing.T) {
		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR, 0777)
		require.Nil(t, err)
		defer f.Close()

		data2 := make([]byte, 0)
		err = f.Decode(&data2)
		require.Nil(t, err)
	})

	t.Run("succeed-5", func(t *testing.T) {
		f, err := OpenFile("./tests/data/test.txt", os.O_CREATE|os.O_RDWR, 0777)
		require.Nil(t, err)
		defer f.Close()

		line, _ := f.ReadLine()
		// require.Nil(t, err)
		fmt.Println(line)
	})
}
