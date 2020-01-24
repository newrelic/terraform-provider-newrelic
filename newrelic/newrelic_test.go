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
	nr, err = New(ConfigAPIKey(testAPIkey), ConfigHTTPTransport(transport))

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

func TestNew_optionInfrastructureBaseURL(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigInfrastructureBaseURL(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("infrastructure base URL can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigInfrastructureBaseURL("http://localhost/"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionSyntheticsBaseURL(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigSyntheticsBaseURL(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("synthetics base URL can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigSyntheticsBaseURL("http://localhost/"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}
func TestNew_optionNerdGraphBaseURL(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigNerdGraphBaseURL(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("nerdgraph base URL can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigNerdGraphBaseURL("http://localhost/"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionLogLevel(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigLogLevel(""))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("log level can not be empty"), err)

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigLogLevel("debug"))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

func TestNew_optionLogJSON(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigLogJSON(true))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}

type TestLogger struct{}

func (t *TestLogger) Error(s string, a ...interface{}) {}
func (t *TestLogger) Info(s string, a ...interface{})  {}
func (t *TestLogger) Debug(s string, a ...interface{}) {}
func (t *TestLogger) Warn(s string, a ...interface{})  {}

func TestNew_optionLogger(t *testing.T) {
	t.Parallel()

	nr, err := New(ConfigAPIKey(testAPIkey), ConfigLogger(nil))
	assert.Nil(t, nr)
	assert.Error(t, errors.New("logger can not be nil"), err)

	testLogger := TestLogger{}

	nr, err = New(ConfigAPIKey(testAPIkey), ConfigLogger(&testLogger))

	assert.NotNil(t, nr)
	assert.NoError(t, err)
}
