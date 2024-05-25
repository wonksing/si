package sifile

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
}
