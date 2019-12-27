// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListDeploymentsResponseJSON = `{
		"deployments": [
			{
				"id": 123,
				"revision": "master",
				"changelog": "v0.0.1",
				"description": "testing",
				"user": "foo",
				"timestamp": "2019-12-27T19:13:23+00:00",
				"links": {
					"application": 111
				}
			}
		]
	}`
)

func TestListDeployments(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testListDeploymentsResponseJSON, http.StatusOK)

	expected := []*Deployment{
		{
			ID:          123,
			Revision:    "master",
			Changelog:   "v0.0.1",
			Description: "testing",
			User:        "foo",
			Timestamp:   "2019-12-27T19:13:23+00:00",
			Links: DeploymentLinks{
				Application: 111,
			},
		},
	}

	actual, err := apm.ListDeployments(123)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
