package newrelic

// resource_newrelic_teams_hierarchy_level manages a TEAMS_HIERARCHY_LEVEL
// entity. These entities are created automatically by NGEP when parent-child
// relationships are established between Team entities via the parentId field.
// From the API's perspective the only user-controllable action is renaming them;
// Terraform therefore supports import + update but not create or delete.
//
// Typical workflow:
//   1. Set parentId on newrelic_team resources to build the hierarchy.
//   2. Discover the resulting TEAMS_HIERARCHY_LEVEL entity IDs via
//      entitySearch (type = 'TEAMS_HIERARCHY_LEVEL') or the Teams UI.
//   3. terraform import newrelic_teams_hierarchy_level.example <entity_id>
//   4. Use this resource to manage the level's name going forward.

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

func resourceNewRelicTeamsHierarchyLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicTeamsHierarchyLevelCreate,
		ReadContext:   resourceNewRelicTeamsHierarchyLevelRead,
		UpdateContext: resourceNewRelicTeamsHierarchyLevelUpdate,
		DeleteContext: resourceNewRelicTeamsHierarchyLevelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The display name for this hierarchy level (e.g. 'Department', 'Squad').",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Tags in 'key:value1,value2' format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NGEP entity type. Always 'TEAMS_HIERARCHY_LEVEL'. Read-only.",
			},
		},
	}
}

// Create is not supported — hierarchy level entities are created automatically
// by NGEP when parent-child Team relationships are established. Use
// terraform import to bring an existing level under Terraform management.
func resourceNewRelicTeamsHierarchyLevelCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "newrelic_teams_hierarchy_level does not support create",
		Detail: "Teams hierarchy level entities are created automatically by the New Relic " +
			"platform when parent-child relationships are established between Team resources " +
			"via the parentId field. Use 'terraform import newrelic_teams_hierarchy_level.<name> <entity_id>' " +
			"to bring an existing level under Terraform management.",
	}}
}

func resourceNewRelicTeamsHierarchyLevelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	entityIface, err := client.Scorecards.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) || isNGEPGhostNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if entityIface == nil {
		d.SetId("")
		return nil
	}

	level, ok := (*entityIface).(*scorecards.EntityManagementTeamsHierarchyLevelEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamsHierarchyLevelEntity", d.Id())
	}

	_ = d.Set("name", level.Name)
	_ = d.Set("type", level.Type)
	_ = d.Set("tags", flattenNGEPTags(level.Tags))

	return nil
}

func resourceNewRelicTeamsHierarchyLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating NGEP teams hierarchy level %s", d.Id())

	upd := scorecards.EntityManagementTeamsHierarchyLevelEntityUpdateInput{}
	if d.HasChange("name") {
		upd.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		upd.Tags = expandNGEPTags(d.Get("tags").(*schema.Set).List())
	}

	if _, err := client.Scorecards.EntityManagementUpdateTeamsHierarchyLevel(d.Id(), upd); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicTeamsHierarchyLevelRead(ctx, d, meta)
}

// Delete is not supported — hierarchy levels are managed by the platform.
// Removing the Terraform resource only drops it from state; the entity persists.
func resourceNewRelicTeamsHierarchyLevelDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[INFO] Removing hierarchy level %s from Terraform state (entity not deleted — managed by NGEP)", d.Id())
	d.SetId("")
	return nil
}
