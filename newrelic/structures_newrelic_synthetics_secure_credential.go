package newrelic

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func expandSyntheticsSecureCredential(d *schema.ResourceData) *synthetics.SecureCredential {
	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	sc := synthetics.SecureCredential{
		Key:         key,
		Value:       d.Get("value").(string),
		Description: d.Get("description").(string),
	}

	return &sc
}

func flattenSyntheticsSecureCredential(sc *synthetics.SecureCredential, d *schema.ResourceData) error {
	d.Set("key", sc.Key)
	d.Set("description", sc.Description)

	createdAt := time.Time(*sc.CreatedAt).Format(time.RFC3339)
	d.Set("created_at", createdAt)

	lastUpdated := time.Time(*sc.CreatedAt).Format(time.RFC3339)
	d.Set("last_updated", lastUpdated)

	return nil
}
