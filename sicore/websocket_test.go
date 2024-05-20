package sicore

import (
	"testing"

	"github.com/wonksing/si/v2/siutils"
)

func Test_NopHub(t *testing.T) {
	nh := NopHub{}
	siutils.AssertNilFail(t, nh.Add(nil))
	siutils.AssertNilFail(t, nh.Remove(nil))
}
