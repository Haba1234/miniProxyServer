package proxyserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_RedirectReq(t *testing.T) {
	t.Parallel()
	s := NewService()

	cases := []struct {
		name         string
		method       string
		url          string
		body         io.Reader
		responseCode int
	}{
		{
			"Method Not Allowed",
			http.MethodGet,
			"http://localhost:8080/proxy",
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			"Unmarshal error",
			http.MethodPost,
			"http://localhost:8080/proxy",
			bytes.NewBufferString(`{"method": "GET", "url": "https://google.com","Headers": {"Authorization": "Basic bG9naW46cGFzc3dvcmQ="}`),
			http.StatusBadRequest,
		},
		{
			"Method wrong",
			http.MethodPost,
			"http://localhost:8080/proxy",
			bytes.NewBufferString(`{"method": "GE", "url": "https://google.com","Headers": {"Authorization": "Basic bG9naW46cGFzc3dvcmQ="}}`),
			http.StatusBadRequest,
		},
		{
			"URL wrong",
			http.MethodPost,
			"http://localhost:8080/proxy",
			bytes.NewBufferString(`{"method": "GET", "url": "google.com","Headers": {"Authorization": "Basic bG9naW46cGFzc3dvcmQ="}}`),
			http.StatusBadRequest,
		},
		{
			"Test OK",
			http.MethodPost,
			"http://localhost:8080/proxy",
			bytes.NewBufferString(`{"method": "GET", "url": "https://google.com","Headers": {"Authorization": "Basic bG9naW46cGFzc3dvcmQ="}}`),
			http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.url, c.body)
			w := httptest.NewRecorder()
			s.RedirectReq(w, r)
			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, c.responseCode, resp.StatusCode)
		})
	}
}
