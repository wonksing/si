package sihttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		e := Error{}
		require.Equal(t, "status: unknown", e.Error())
	})

	t.Run("succeed-2", func(t *testing.T) {
		e := Error{
			Response: &http.Response{
				Status: http.StatusText(http.StatusBadRequest),
			},
			Body: []byte("bad request"),
		}
		require.Equal(t, "status: "+http.StatusText(http.StatusBadRequest)+", body: bad request", e.Error())
	})
}

func TestError_GetStatusCode(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		e := Error{}
		require.Equal(t, http.StatusInternalServerError, e.GetStatusCode(http.StatusInternalServerError))
	})

	t.Run("succeed-2", func(t *testing.T) {
		e := Error{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     http.StatusText(http.StatusBadRequest),
			},
			Body: []byte("bad request"),
		}
		require.Equal(t, http.StatusBadRequest, e.GetStatusCode(http.StatusBadRequest))
	})
}

func TestError_GetStatus(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		e := Error{}
		require.Equal(t, http.StatusText(http.StatusInternalServerError), e.GetStatus(http.StatusInternalServerError))
	})

	t.Run("succeed-2", func(t *testing.T) {
		e := Error{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     http.StatusText(http.StatusBadRequest),
			},
			Body: []byte("bad request"),
		}
		require.Equal(t, http.StatusText(http.StatusBadRequest), e.GetStatus(http.StatusBadRequest))
	})
}
