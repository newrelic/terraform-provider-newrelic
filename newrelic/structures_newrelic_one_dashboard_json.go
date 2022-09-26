package newrelic

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
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

func flattenDashboardJsonEntity(dashboard *entities.DashboardEntity, d *schema.ResourceData) error {
	_ = d.Set("account_id", dashboard.AccountID)
	_ = d.Set("guid", dashboard.GUID)
	_ = d.Set("permalink", dashboard.Permalink)

	json, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	_ = d.Set("hash_remote", hashString(json))

	return nil
}

func hashString(content []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(content))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
