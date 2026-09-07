//go:build unit

package newrelic

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
)

// newMockEntitySearchServer serves a sequence of canned NerdGraph responses, one per
// request, and records the raw request body of each request so tests can assert on
// the GraphQL variables actually sent over the wire (in particular, the cursor).
func newMockEntitySearchServer(t *testing.T, responses []string) (*nr.NewRelic, *[][]byte) {
	var requestBodies [][]byte
	callCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requestBodies = append(requestBodies, body)

		require.Less(t, callCount, len(responses), "unexpected extra request to mock NerdGraph server")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(responses[callCount]))
		require.NoError(t, err)

		callCount++
	}))
	t.Cleanup(ts.Close)

	client, err := nr.New(
		nr.ConfigPersonalAPIKey("test-api-key"),
		nr.ConfigNerdGraphBaseURL(ts.URL),
	)
	require.NoError(t, err)

	return client, &requestBodies
}

func entitySearchPage(entitiesJSON, nextCursor string) string {
	return `{"data":{"actor":{"entityManagement":{"entitySearch":{"entities":[` +
		entitiesJSON + `],"nextCursor":"` + nextCursor + `"}}}}}`
}

func agentConfigEntityJSON(id, name, scopeType, scopeID string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"__typename": "EntityManagementAgentConfigurationEntity",
		"id":         id,
		"name":       name,
		"type":       "AGENT_CONFIGURATION",
		"scope": map[string]string{
			"id":   scopeID,
			"type": scopeType,
		},
	})
	return string(b)
}

type entitySearchRequestBody struct {
	Variables map[string]interface{} `json:"variables"`
}

func TestUnitFindFleetConfigurationsByName_FollowsCursorAcrossPages(t *testing.T) {
	t.Parallel()

	page1 := entitySearchPage(agentConfigEntityJSON("cfg-other", "other-config", "ORGANIZATION", "org-1"), "page-2-cursor")
	page2 := entitySearchPage(agentConfigEntityJSON("cfg-match", "target-config", "ORGANIZATION", "org-1"), "")

	client, requestBodies := newMockEntitySearchServer(t, []string{page1, page2})

	matches, err := findFleetConfigurationsByName(client, "target-config")
	require.NoError(t, err)
	require.Len(t, matches, 1)
	assert.Equal(t, "cfg-match", matches[0].ID)

	require.Len(t, *requestBodies, 2, "expected one request per page")

	var firstReq, secondReq entitySearchRequestBody
	require.NoError(t, json.Unmarshal((*requestBodies)[0], &firstReq))
	require.NoError(t, json.Unmarshal((*requestBodies)[1], &secondReq))

	assert.Nil(t, firstReq.Variables["cursor"], "first request should search from the beginning")
	assert.Equal(t, "page-2-cursor", secondReq.Variables["cursor"], "second request should use the cursor returned by page 1")
}

func TestUnitFindFleetConfigurationsByName_MultipleMatchesAcrossPages(t *testing.T) {
	t.Parallel()

	page1 := entitySearchPage(agentConfigEntityJSON("cfg-1", "shared-name", "ORGANIZATION", "org-1"), "page-2-cursor")
	page2 := entitySearchPage(agentConfigEntityJSON("cfg-2", "shared-name", "ORGANIZATION", "org-1"), "")

	client, _ := newMockEntitySearchServer(t, []string{page1, page2})

	matches, err := findFleetConfigurationsByName(client, "shared-name")
	require.NoError(t, err)
	require.Len(t, matches, 2, "matches with the same name on different pages should all be returned")

	ids := []string{matches[0].ID, matches[1].ID}
	assert.ElementsMatch(t, []string{"cfg-1", "cfg-2"}, ids)
}

func TestUnitFindFleetConfigurationsByName_NoMatch(t *testing.T) {
	t.Parallel()

	page1 := entitySearchPage(agentConfigEntityJSON("cfg-1", "unrelated-config", "ORGANIZATION", "org-1"), "")

	client, _ := newMockEntitySearchServer(t, []string{page1})

	matches, err := findFleetConfigurationsByName(client, "does-not-exist")
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestUnitFindFleetConfigurationsByName_StopsWhenCursorEmpty(t *testing.T) {
	t.Parallel()

	// Only one page is registered; if the function kept paging past an empty
	// nextCursor it would request a second page and the mock server would fail
	// the test via the callCount bounds check in newMockEntitySearchServer.
	page1 := entitySearchPage(agentConfigEntityJSON("cfg-1", "target-config", "ORGANIZATION", "org-1"), "")

	client, requestBodies := newMockEntitySearchServer(t, []string{page1})

	matches, err := findFleetConfigurationsByName(client, "target-config")
	require.NoError(t, err)
	require.Len(t, matches, 1)
	assert.Len(t, *requestBodies, 1)
}
