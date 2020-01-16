// +build unit

package newrelic

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var testAPIkey = "asdf1234"

func TestNew_invalid(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(""))

	assert.Nil(t, nr)
	assert.Error(t, errors.New("apiKey required"), err)
}

func TestNew_basic(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_configOptionError(t *testing.T) {
	t.Parallel()

	badOption := func(cfg *config.Config) error { return errors.New("option with error") }
	nr, err := New(ConfigAPIKey(testAPIkey), badOption)

	assert.Nil(t, nr)
	assert.Error(t, errors.New("option with error"), err)
}

func TestNew_setPersonalAPIKey(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigPersonalAPIKey(testAPIkey))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_setRegion(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigRegion("US"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionTimeout(t *testing.T) {
	t.Parallel()

	timeout := time.Second * 30
	nr, err := New(ConfigAPIKey(testAPIkey), ConfigHTTPTimeout(timeout))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionTransport(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigHTTPTransport(nil))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("HTTP Transport can not be nil"), err)

	transport := http.DefaultTransport
	nr, err = New(ConfigAPIKey(testAPIkey), ConfigHTTPTransport(&transport))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionUserAgent(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigUserAgent(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("user-agent can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigUserAgent("my-user-agent"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionBaseURL(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigBaseURL(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("base URL can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigBaseURL("http://localhost/"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}
