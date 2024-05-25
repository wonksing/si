package siwebsocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NopHub(t *testing.T) {
	nh := NopHub{}
	require.Nil(t, nh.Add(nil))
	require.Nil(t, nh.Remove(nil))
}
