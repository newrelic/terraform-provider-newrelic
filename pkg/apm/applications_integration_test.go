// +build integration

package apm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationApplications(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	a, err := client.ListApplications(nil)
	require.NoError(t, err)

	_, err = client.GetApplication(a[0].ID)
	require.NoError(t, err)

	params := UpdateApplicationParams{
		Name:     a[0].Name,
		Settings: a[0].Settings,
	}

	_, err = client.UpdateApplication(a[0].ID, params)
	require.NoError(t, err)

	n, err := client.GetMetricNames(a[0].ID, MetricNamesParams{})
	require.NoError(t, err)

	_, err = client.GetMetricData(a[0].ID, MetricDataParams{Names: []string{n[0].Name}})
	require.NoError(t, err)
}

func TestIntegrationDeleteApplication(t *testing.T) {
	t.Skip("What does delete mean in the case where we have no create?")
	t.Parallel()

	client := newIntegrationTestClient(t)

	_, err := client.DeleteApplication(0)

	if err != nil {
		t.Fatal(err)
	}
}
