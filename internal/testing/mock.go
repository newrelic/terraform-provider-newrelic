package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// NewMockResponse creates a server to respond to API calls for unit tests
func NewMockResponse(t *testing.T, mockJSONResponse string, statusCode int, contentType string) *httptest.Server {
	if contentType == "" {
		contentType = "application/json"
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))

	return ts
}
