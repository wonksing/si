package sifile_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sifile"
)

func TestListAll(t *testing.T) {
	list, err := sifile.ListDir("./tests")
	require.Nil(t, err)

	for _, f := range list {
		fi, err := f.Info()
		require.Nil(t, err)
		fmt.Println(f.Path, f.IsDir(), fi.Size())
	}
}
