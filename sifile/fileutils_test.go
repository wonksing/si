package sifile_test

import (
	"fmt"
	"testing"

	"github.com/wonksing/si/sifile"
	"github.com/wonksing/si/siutils"
)

func TestListAll(t *testing.T) {
	list, err := sifile.ListDir("./tests")
	siutils.AssertNilFail(t, err)

	for _, f := range list {
		fi, err := f.Info()
		siutils.AssertNilFail(t, err)
		fmt.Println(f.Path, f.IsDir(), fi.Size())
	}
}
