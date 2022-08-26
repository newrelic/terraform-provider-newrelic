//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/ai"

	"github.com/stretchr/testify/assert"
)

func testFlattenNotificationDestinationAuth(t *testing.T, v interface{}, auth ai.AiNotificationsAuth) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "type":
			assert.Equal(t, cv, string(auth.AuthType))
		case "user":
			assert.Equal(t, cv, auth.User)
		}
	}
}
