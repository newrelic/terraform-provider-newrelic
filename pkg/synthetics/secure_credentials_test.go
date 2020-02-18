// +build unit

package synthetics

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testSecureCredentialKey         = "TEST"
	testSecureCredentialDescription = "Test Description"
	testSecureCredential            = &SecureCredential{
		Key:         testSecureCredentialKey,
		Description: testSecureCredentialDescription,
		CreatedAt:   &testTimestamp,
		LastUpdated: &testTimestamp,
	}
	testGetSecureCredentialsJson = fmt.Sprintf(`
	{
		"secureCredentials": [
			{
				"key": "%[1]s",
				"description": "%[2]s",
				"createdAt": "2019-11-27T19:11:05.076+0000",
				"lastUpdated": "2019-11-27T19:11:05.076+0000"
			}, {
				"key": "myKey2",
				"description": "Description of myKey2",
				"createdAt": "2019-11-27T19:11:05.076+0000",
				"lastUpdated": "2019-11-27T19:11:05.076+0000"
			}
		],
		"count": 2
	}
	`, testSecureCredentialKey, testSecureCredentialDescription, testTimestamp)
)

func TestGetSecureCredentials(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, testGetSecureCredentialsJson, http.StatusOK)

	r, err := synthetics.GetSecureCredentials()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(r))
	assert.Equal(t, r[0], testSecureCredential)
}

func TestAddSecureCredential(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, "", http.StatusOK)

	err := synthetics.AddMonitorLabel(testSecureCredentialKey, "test", "test")
	assert.NoError(t, err)
}

func TestDeleteSecureCredenti(t *testing.T) {
	t.Parallel()
	synthetics := newMockResponse(t, "", http.StatusOK)

	err := synthetics.DeleteMonitorLabel(testSecureCredentialKey, "test", "test")
	assert.NoError(t, err)
}
