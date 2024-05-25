package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultEofChecker_Check(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		c := defaultEofChecker{}

		inputErr := errors.New("just error")
		res, err := c.Check(nil, inputErr)
		require.NotNil(t, err)
		require.Equal(t, inputErr, err)
		require.EqualValues(t, false, res)
	})
}
