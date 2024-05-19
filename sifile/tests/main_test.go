package sifile_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
)

func setup() error {
	return nil
}

func shutdown() {
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
}
