// +build integration

package dashboards

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDashboards(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	dashboards := New(config.Config{
		APIKey: apiKey,
	})

	d := Dashboard{
		Metadata: DashboardMetadata{
			Version: 1,
		},
		Title:      "newrelic-client-go-test",
		Visibility: Visibility.Owner,
		Editable:   Editable.Owner,
	}

	// Test: Create
	created, err := dashboards.CreateDashboard(d)
	if err != nil {
		t.Fatalf("CreateDashboards error: %s", err)
	}

	assert.NotNil(t, created)

	// Test: List
	params := ListDashboardsParams{
		Title: "newrelic-client-go",
	}
	multiple, err := dashboards.ListDashboards(&params)
	if err != nil {
		t.Fatalf("ListDashboards error: %s", err)
	}

	assert.NotNil(t, multiple)

	// Test: Get
	single, err := dashboards.GetDashboard(multiple[0].ID)
	if err != nil {
		t.Fatalf("GetDashboards error: %s", err)
	}

	assert.NotNil(t, single)

	// Test: Update
	single.Title = "updated"
	updated, err := dashboards.UpdateDashboard(*single)
	if err != nil {
		t.Fatalf("UpdateDashboards error: %s", err)
	}

	assert.NotNil(t, updated)

	// Test: Delete
	deleted, err := dashboards.DeleteDashboard(updated.ID)
	if err != nil {
		t.Fatalf("DeleteDashboards error: %s", err)
	}

	assert.NotNil(t, deleted)
}
