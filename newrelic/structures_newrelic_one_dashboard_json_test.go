package newrelic

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
)

func TestSetDataAccounts(t *testing.T) {
	const fixtureInName = "dashboard_data_account.json"
	const fixtureOutName = "dashboard_data_account_expected.json"

	inputFile, err := os.ReadFile(filepath.Join("testdata", fixtureInName))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	dash := dashboards.DashboardInput{}
	if err = json.Unmarshal(inputFile, &dash); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	dataAccounts := []any{456, 789}
	if err = setDataAccounts(&dash, dataAccounts); err != nil {
		t.Fatalf("setDataAccounts: %v", err)
	}

	var expectedRaw []byte
	expectedRaw, err = os.ReadFile(filepath.Join("testdata", fixtureOutName))
	if err != nil {
		t.Fatalf("read reference fixture: %v", err)
	}
	var expected dashboards.DashboardInput
	if err = json.Unmarshal(expectedRaw, &expected); err != nil {
		t.Fatalf("unmarshal reference fixture: %v", err)
	}

	if diff := cmp.Diff(expected, dash); diff != "" {
		t.Errorf("expected differs from actual: %s", diff)
	}
}
