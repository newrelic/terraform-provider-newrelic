package newrelic

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicCompositeServiceLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCompositeServiceLevelCreate,
		ReadContext:   resourceNewRelicCompositeServiceLevelRead,
		UpdateContext: resourceNewRelicCompositeServiceLevelUpdate,
		DeleteContext: resourceNewRelicCompositeServiceLevelDelete,
		Importer: &schema.ResourceImporter{
			// Composite's own GUID is the full resource ID — no entity-scoped
			// read dependency like the plain service_level resource has, so no
			// custom identifier parsing is needed here.
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true, // immutable — the update mutation has no accountId field at all
				Description:  "The New Relic account ID where the composite will reside.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the composite.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the composite.",
			},
			"component": {
				// TypeSet, not TypeList: the server doesn't guarantee it returns
				// components in submission order (confirmed via acceptance
				// testing — order came back swapped from what was sent), and
				// composites are semantically a set of components anyway (no
				// validation rule cares about order). A List would require
				// reordering server responses to match config on every Read,
				// which breaks down specifically on import (no prior config to
				// anchor against). Set's only real downside — silently
				// collapsing two components with an *identical* guid AND
				// weight — is narrower and caught by validateCompositeComponents
				// in the one sub-case that actually matters (different
				// weights, or no weights at all, don't collapse).
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    2,
				MaxItems:    20,
				Description: "The component service levels that make up this composite. Must be 2-20 components. Reference other resources' `.guid` output rather than a hardcoded GUID string, so Terraform's dependency graph deletes resources in the right order.",
				Elem:        compositeComponentSchema(),
			},
			"objective": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The objective definition for the composite.",
				Elem:        objectiveSchema(), // reused from resource_newrelic_service_level.go — same package, same shape, server enforces exactly 1 objective for composites too
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the composite.",
			},
		},
	}
}

func compositeComponentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique entity identifier of the component service level.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"weight": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The weight of this service level in the composite calculation, expressed as a percentage (0-100). All-or-nothing: if any component has a weight, every component must have one, and they must sum to exactly 100. If none are given, the server distributes weight equally.",
			},
		},
	}
}

// validateCompositeComponents checks the two component-level rules the
// server itself enforces (slm-api's FetcherValidator: no duplicate GUIDs,
// weights are all-or-nothing summing to exactly 100) so users get a fast,
// local error instead of a round-trip to the API. Component count (2-20) is
// already enforced by the schema's MinItems/MaxItems. Cross-account
// authorization, nesting depth/cycle checks, and delete-protection are
// intentionally NOT checked here — they require live server state the
// provider doesn't have at plan time.
func validateCompositeComponents(components []interface{}) error {
	seenGUIDs := make(map[string]bool, len(components))
	weightedCount := 0
	var weightSum float64

	for i, c := range components {
		cfg := c.(map[string]interface{})
		guid := cfg["guid"].(string)

		if seenGUIDs[guid] {
			return fmt.Errorf("component[%d]: duplicate guid %q — no two components may reference the same service level", i, guid)
		}
		seenGUIDs[guid] = true

		if weight, ok := cfg["weight"].(float64); ok && weight != 0 {
			weightedCount++
			weightSum += weight
		}
	}

	if weightedCount > 0 && weightedCount != len(components) {
		return fmt.Errorf("weights are all-or-nothing: if any component has a weight, every component must have one (got %d of %d components with a weight)", weightedCount, len(components))
	}

	if weightedCount > 0 && math.Abs(weightSum-100) > 1e-6 {
		return fmt.Errorf("component weights must sum to exactly 100, got %g", weightSum)
	}

	return nil
}

func resourceNewRelicCompositeServiceLevelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	components := d.Get("component").(*schema.Set).List()
	if err := validateCompositeComponents(components); err != nil {
		return diag.FromErr(err)
	}

	createInput := expandCompositeServiceLevelCreateInput(d)

	log.Printf("[INFO] Creating New Relic One Composite Service Level %s", createInput.Name)

	created, err := client.ServiceLevel.ServiceLevelCompositeCreateWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(created.GUID))

	return diag.FromErr(flattenCompositeServiceLevelIndicator(*created, d))
}

func resourceNewRelicCompositeServiceLevelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	composites, err := client.ServiceLevel.GetCompositesWithContext(ctx, []common.EntityGUID{common.EntityGUID(d.Id())})
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	for _, composite := range *composites {
		if string(composite.GUID) == d.Id() {
			return diag.FromErr(flattenCompositeServiceLevelIndicator(composite, d))
		}
	}

	d.SetId("")
	return nil
}

func resourceNewRelicCompositeServiceLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	components := d.Get("component").(*schema.Set).List()
	if err := validateCompositeComponents(components); err != nil {
		return diag.FromErr(err)
	}

	updateInput := expandCompositeServiceLevelUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One Composite Service Level %s", d.Id())

	_, err := client.ServiceLevel.ServiceLevelCompositeUpdateWithContext(ctx, updateInput, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicCompositeServiceLevelRead(ctx, d, meta)
}

func resourceNewRelicCompositeServiceLevelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One Composite Service Level %s", d.Id())

	if _, err := client.ServiceLevel.ServiceLevelCompositeDeleteWithContext(ctx, common.EntityGUID(d.Id())); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
