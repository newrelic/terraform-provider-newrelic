// +build integration

package synthetics

import (
	"fmt"
	"os"
	"strings"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testIntegrationSecureCredentialValue       = "test value"
	testIntegrationSecureCredentialDescription = "test description"
	testIntegrationSecureCredential            = &SecureCredential{
		Value:       testIntegrationSecureCredentialValue,
		Description: testIntegrationSecureCredentialDescription,
	}
)

func TestIntegrationSecureCredentials(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	synthetics := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	// Setup
	rand := strings.ToUpper(nr.RandSeq(5))
	key := fmt.Sprintf("TEST_SYNTHETICS_SECURE_CREDENTIAL_%s", rand)
	testIntegrationSecureCredential.Key = key
	creds, err := synthetics.GetSecureCredentials()
	require.NoError(t, err)

	originalCount := len(creds)

	// Test: Add
	c, err := synthetics.AddSecureCredential(key, "asdf", testIntegrationSecureCredentialDescription)
	require.NoError(t, err)

	// Test: Get
	c, err = synthetics.GetSecureCredential(key)
	require.NoError(t, err)
	assert.Equal(t, (*c).Key, key)
	assert.Equal(t, (*c).Description, testIntegrationSecureCredentialDescription)

	// Test: Get (Multiple)
	creds, err = synthetics.GetSecureCredentials()
	require.NoError(t, err)
	assert.Equal(t, originalCount+1, len(creds))

	// Test: Update
	c, err = synthetics.UpdateSecureCredential(c.Key, testIntegrationSecureCredentialValue, "new test value")
	require.NoError(t, err)
	assert.Equal(t, "new test value", (*c).Description)

	// Test: Delete
	err = synthetics.DeleteSecureCredential(key)
	require.NoError(t, err)

	// Verify Delete
	_, err = synthetics.GetSecureCredential(key)
	require.Error(t, err)
}
