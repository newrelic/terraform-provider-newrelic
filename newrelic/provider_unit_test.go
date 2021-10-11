//go:build unit
// +build unit

package newrelic

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfig(t *testing.T) {
	c := ProviderConfig{
		PersonalAPIKey: "abc123",
		AccountID:      123,
	}

	hasNerdGraphCreds := c.hasNerdGraphCredentials()

	if !hasNerdGraphCreds {
		t.Error("hasNerdGraphCreds should be true")
	}
}
