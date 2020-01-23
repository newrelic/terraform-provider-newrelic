package entities

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSetNerdGraphBaseURL(t *testing.T) {
	a := New(config.Config{
		NerdGraphBaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.client.Client.Config.BaseURL)
}
