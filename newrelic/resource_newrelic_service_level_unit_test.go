//go:build unit

package newrelic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
)

const (
	testServiceLevelAccountID  = 1234567
	testServiceLevelSliID      = "2910461"
	testServiceLevelEntityGUID = "MTIzNDU2N3xFWFR8U0VSVklDRXwx"
)

// newMockServiceLevelServer answers the indicators query with the canned response registered for
// the entity GUID in the request variables, and records which GUIDs were queried.
func newMockServiceLevelServer(t *testing.T, responses map[string]string) (*nr.NewRelic, *[]string) {
	var queried []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req struct {
			Variables struct {
				EntityGUID string `json:"entityGUID"`
			} `json:"variables"`
		}
		require.NoError(t, json.Unmarshal(body, &req))
		queried = append(queried, req.Variables.EntityGUID)

		resp, ok := responses[req.Variables.EntityGUID]
		require.True(t, ok, "unexpected query for entity %s", req.Variables.EntityGUID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(resp))
		require.NoError(t, err)
	}))
	t.Cleanup(ts.Close)

	client, err := nr.New(
		nr.ConfigPersonalAPIKey("test-api-key"),
		nr.ConfigNerdGraphBaseURL(ts.URL),
	)
	require.NoError(t, err)

	return client, &queried
}

func serviceLevelIndicatorsResponse(ids ...string) string {
	indicators := make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		indicators = append(indicators, map[string]interface{}{
			"id":         id,
			"name":       "sli-" + id,
			"entityGuid": testServiceLevelEntityGUID,
			"events": map[string]interface{}{
				"account":     map[string]interface{}{"id": testServiceLevelAccountID},
				"validEvents": map[string]interface{}{"from": "Transaction", "where": "appName = 'Example'"},
				"goodEvents":  map[string]interface{}{"from": "Transaction", "where": "appName = 'Example' AND duration < 0.1"},
			},
			"objectives": []map[string]interface{}{{
				"target":     99,
				"timeWindow": map[string]interface{}{"rolling": map[string]interface{}{"count": 7, "unit": "DAY"}},
			}},
		})
	}
	b, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entity": map[string]interface{}{"serviceLevel": map[string]interface{}{"indicators": indicators}},
			},
		},
	})
	return string(b)
}

// The response NerdGraph gives for an entity that does not exist (or is not indexed yet).
const serviceLevelEntityNotFoundResponse = `{"data":{"actor":{"entity":null}}}`

func readServiceLevel(t *testing.T, responses map[string]string) (*schema.ResourceData, []string) {
	identifier := &serviceLevelIdentifier{
		AccountID:  testServiceLevelAccountID,
		ID:         testServiceLevelSliID,
		EntityGUID: testServiceLevelEntityGUID,
	}
	client, queried := newMockServiceLevelServer(t, responses)

	d := schema.TestResourceDataRaw(t, resourceNewRelicServiceLevel().Schema, map[string]interface{}{})
	d.SetId(identifier.String())

	diags := resourceNewRelicServiceLevelRead(context.Background(), d, &ProviderConfig{NewClient: client})
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	return d, *queried
}

func testServiceLevelSliGUID() string {
	return getSliGUID(&serviceLevelIdentifier{AccountID: testServiceLevelAccountID, ID: testServiceLevelSliID})
}

func TestUnitServiceLevelRead_FoundBySliGUID(t *testing.T) {
	sliGUID := testServiceLevelSliGUID()

	d, queried := readServiceLevel(t, map[string]string{
		sliGUID: serviceLevelIndicatorsResponse(testServiceLevelSliID),
	})

	require.NotEmpty(t, d.Id())
	require.Equal(t, "sli-"+testServiceLevelSliID, d.Get("name"))
	require.Equal(t, []string{sliGUID}, queried)
}

func TestUnitServiceLevelRead_SliEntityNotIndexedYet(t *testing.T) {
	sliGUID := testServiceLevelSliGUID()

	d, queried := readServiceLevel(t, map[string]string{
		sliGUID:                    serviceLevelEntityNotFoundResponse,
		testServiceLevelEntityGUID: serviceLevelIndicatorsResponse("1111111", testServiceLevelSliID),
	})

	require.NotEmpty(t, d.Id(), "a just-created SLI must not be dropped from the state")
	require.Equal(t, "sli-"+testServiceLevelSliID, d.Get("name"))
	require.Equal(t, sliGUID, d.Get("sli_guid"))
	require.Equal(t, []string{sliGUID, testServiceLevelEntityGUID}, queried)
}

func TestUnitServiceLevelRead_DeletedSli(t *testing.T) {
	sliGUID := testServiceLevelSliGUID()

	d, _ := readServiceLevel(t, map[string]string{
		sliGUID:                    serviceLevelEntityNotFoundResponse,
		testServiceLevelEntityGUID: serviceLevelIndicatorsResponse("1111111"),
	})

	require.Empty(t, d.Id())
}

func TestUnitServiceLevelRead_DeletedSliAndEntity(t *testing.T) {
	sliGUID := testServiceLevelSliGUID()

	d, _ := readServiceLevel(t, map[string]string{
		sliGUID:                    serviceLevelEntityNotFoundResponse,
		testServiceLevelEntityGUID: serviceLevelEntityNotFoundResponse,
	})

	require.Empty(t, d.Id())
}
