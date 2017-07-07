package api

import (
	"net/http"
	"testing"
)

func TestQueryApplications_Basic(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
      {
        "applications": [
          {
            "id": 123,
            "name": "foo"
          },
          {
            "id": 456,
            "name": "bar"
          }
        ]
      }
    `))
	}))

	apps, err := c.queryApplications(applicationsFilters{})
	if err != nil {
		t.Log(err)
		t.Fatal("queryApplications error")
	}

	if len(apps) == 0 {
		t.Fatal("No applications found")
	}
}
