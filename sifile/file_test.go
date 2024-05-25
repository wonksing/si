package sifile

import (
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
