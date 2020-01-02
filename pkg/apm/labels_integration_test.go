// +build integration

package apm

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationLabels(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Setup
	appFilterName := "tf_" // Filter for integration test applications
	appQueryParams := ListApplicationsParams{
		Name: &appFilterName,
	}
	applications, err := client.ListApplications(&appQueryParams)

	if err != nil {
		t.Skipf("Skipped Labels integration test due error fetching applications: %s", err)
		return
	}

	if len(applications) == 0 {
		t.Skipf("Skipped Labels integration test due to no applications being found for filter by name: %s", appFilterName)
		return
	}

	application := applications[0]
	labelName := nr.RandSeq(5)
	newLabel := Label{
		Category: "Project",
		Name:     fmt.Sprintf("label-%s", labelName),
		Links: LabelLinks{
			Applications: []int{application.ID},
			Servers:      []int{},
		},
	}

	t.Logf("applying label %s linked to application %s", labelName, application.Name)

	// Test: Create
	label, err := client.CreateLabel(newLabel)

	require.NoError(t, err)
	require.NotNil(t, label)

	// Test: List
	listResult, err := client.ListLabels()

	require.NoError(t, err)
	require.Greater(t, len(listResult), 0)

	// Test: Get
	getResult, err := client.GetLabel(label.Key)

	require.NoError(t, err)
	require.NotNil(t, getResult)

	// Test: Delete
	deleteResult, err := client.DeleteLabel(label.Key)

	require.NoError(t, err)
	require.NotNil(t, deleteResult)
}
