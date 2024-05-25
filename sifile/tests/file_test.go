package sifile_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sifile"
)

func TestFile_ReadFrom(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	data := "test data to write.\n"
	dataReader := strings.NewReader(data)
	fileName := "./data/TestFile_ReadFrom.txt"

	var fileMode os.FileMode
	fi, err := os.Stat(fileName)
	if err != nil {
		fileMode = 0755
	} else {
		fileMode = fi.Mode()
	}
	f, err := sifile.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, fileMode)
	require.Nil(t, err)
	defer f.Close()

	n, err := f.ReadFrom(dataReader)
	require.Nil(t, err)

	fmt.Println(n)
}
