package siftp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/experimental/siftp"
)

func TestRead(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	c := siftp.NewClient("", "", "")
	res, err := c.ReadFile("Version.ini")
	require.Nil(t, err)

	fmt.Println(string(res))
}

func TestWriteFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	c := siftp.NewClient("", "", "")
	err := c.WriteFile("test.txt", []byte("test upload"))
	require.Nil(t, err)

}
