package newrelic

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
)

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardJsonInput(d *schema.ResourceData, meta interface{}) (*dashboards.DashboardInput, error) {
	dash := dashboards.DashboardInput{}
	err := json.Unmarshal([]byte(d.Get("json").(string)), &dash)
	if err != nil {
		return nil, err
	}

	return &dash, nil
}

func hashString(content []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(content))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
