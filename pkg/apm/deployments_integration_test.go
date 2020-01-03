// +build integration

package apm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployments(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Setup
	appFilterName := "tf_" // Filter for integration test applications
	appQueryParams := ListApplicationsParams{
		Name: appFilterName,
	}
	applications, err := client.ListApplications(&appQueryParams)

	if err != nil {
		t.Skipf("Skipped Deployments integration test due error fetching applications: %s", err)
		return
	}

	if len(applications) == 0 {
		t.Skipf("Skipped Deployments integration test due to no applications being found for filter by name: %s", appFilterName)
		return
	}

	application := applications[0]
	newDeployment := Deployment{
		Revision:    "master",
		Changelog:   "v0.0.1",
		Description: "testing",
		User:        "foo",
	}

	// Test: Create
	deployment, err := client.CreateDeployment(application.ID, newDeployment)

	require.NoError(t, err)
	require.NotNil(t, deployment)

	// Test: List
	listResult, err := client.ListDeployments(application.ID)

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Delete
	result, err := client.DeleteDeployment(application.ID, deployment.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
}
