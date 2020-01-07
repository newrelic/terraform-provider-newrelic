// +build unit

package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	testBaseURL := "https://www.mocky.io"
	testTimeout := time.Second * 5
	testTransport := http.DefaultTransport

	c := NewClient(config.Config{
		APIKey:        testAPIKey,
		BaseURL:       testBaseURL,
		UserAgent:     testUserAgent,
		Timeout:       &testTimeout,
		HTTPTransport: &testTransport,
	})

	assert.Equal(t, &testTimeout, c.Config.Timeout)
	assert.Equal(t, testBaseURL, c.Config.BaseURL)
	assert.Same(t, &testTransport, c.Config.HTTPTransport)
}

func TestConfigDefaults(t *testing.T) {
	t.Parallel()
	c := NewClient(config.Config{
		APIKey: testAPIKey,
	})

	assert.Equal(t, config.DefaultBaseURLs[c.Config.Region], c.Config.BaseURL)
	assert.Contains(t, c.Config.UserAgent, "newrelic/newrelic-client-go/")
}

func TestDefaultErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"title":"error message"}}`))
	}))

	_, err := c.Get("/path", nil, nil)

	assert.Equal(t, err.(*errors.ErrorUnexpectedStatusCode).Err, "error message")
}

type CustomErrorResponse struct {
	CustomError string `json:"custom"`
}

func (c *CustomErrorResponse) Error() string {
	return c.CustomError
}

func TestCustomErrorValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"custom":"error message"}`))
	}))

	c.SetErrorValue(&CustomErrorResponse{})

	_, err := c.Get("/path", nil, nil)

	assert.Equal(t, err.(*errors.ErrorUnexpectedStatusCode).Err, "error message")
}

type CustomResponseValue struct {
	Custom string `json:"custom"`
}

func TestResponseValue(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"custom":"custom response string"}`))
	}))

	v := &CustomResponseValue{}
	_, err := c.Get("/path", nil, v)

	assert.NoError(t, err)
	assert.Equal(t, &CustomResponseValue{Custom: "custom response string"}, v)
}

func TestQueryParams(t *testing.T) {
	t.Parallel()
	queryParams := struct {
		A int `url:"a,omitempty"`
		B int `url:"b,omitempty"`
	}{
		A: 1,
		B: 2,
	}

	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		a := r.URL.Query().Get("a")
		assert.Equal(t, "1", a)

		b := r.URL.Query().Get("b")
		assert.Equal(t, "2", b)
	}))

	_, _ = c.Get("/path", &queryParams, nil)
}

type TestRequestBody struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestRequestBodyMarshal(t *testing.T) {
	t.Parallel()
	expected := TestRequestBody{
		A: "1",
		B: "2",
	}

	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		actual := &TestRequestBody{}
		err := json.NewDecoder(r.Body).Decode(&actual)

		assert.NoError(t, err)
		assert.Equal(t, &expected, actual)
	}))

	_, _ = c.Post("/path", nil, expected, nil)
}

type TestInvalidRequestBody struct {
	Channel chan int `json:"a"`
}

func TestRequestBodyMarshalError(t *testing.T) {
	t.Parallel()
	b := TestInvalidRequestBody{
		Channel: make(chan int),
	}

	c := NewTestAPIClient(nil)

	_, err := c.Post("/path", nil, b, nil)
	assert.Error(t, err)
}

func TestUrlParseError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(nil)

	_, err := c.Get("\\", nil, nil)
	assert.Error(t, err)
}

func TestPathOnlyUrl(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https://www.mocky.io/v2/path")
	}))

	c.Config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("/path", nil, nil)
}

func TestHostAndPathUrl(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, r.URL, "https:/www.httpbin.org/path")
	}))

	c.Config.BaseURL = "https://www.mocky.io/v2"

	_, _ = c.Get("https://www.httpbin.org/path", nil, nil)
}

type TestInvalidReponseBody struct {
	Channel chan int `json:"channel"`
}

func TestResponseUnmarshalError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"channel": "test"}`))
	}))

	_, err := c.Get("/path", nil, &TestInvalidReponseBody{})

	assert.Error(t, err)
}

func TestHeaders(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, testAPIKey, r.Header.Get("x-api-key"))
		assert.Equal(t, testUserAgent, r.Header.Get("user-agent"))
	}))

	_, err := c.Get("/path", nil, nil)

	assert.Nil(t, err)
}

func TestErrNotFound(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	_, err := c.Get("/path", nil, nil)

	assert.IsType(t, &errors.ErrorNotFound{}, err)
}

func TestInternalServerError(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := c.Get("/path", nil, nil)

	assert.IsType(t, &errors.ErrorUnexpectedStatusCode{}, err)
}

func TestPost(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Post("/path", &struct{}{}, &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}

func TestPut(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Put("/path", &struct{}{}, &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	c := NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{}`))
	}))

	_, err := c.Delete("/path", &struct{}{}, &struct{}{})

	assert.NoError(t, err)
}
