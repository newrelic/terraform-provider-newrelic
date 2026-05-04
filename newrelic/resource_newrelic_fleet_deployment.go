package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleetDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetDeploymentCreate,
		ReadContext:   resourceNewRelicFleetDeploymentRead,
		UpdateContext: resourceNewRelicFleetDeploymentUpdate,
		DeleteContext: resourceNewRelicFleetDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceNewRelicFleetDeploymentCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the fleet this deployment belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the deployment.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the deployment.",
			},
			"agent": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "One or more agent blocks. Each agent type may appear at most once per deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"NRInfra",
								"NRDOT",
								"FluentBit",
								"NRPrometheusAgent",
							}, false),
							Description: "The agent type. Allowed values: NRInfra, NRDOT, FluentBit, NRPrometheusAgent.",
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The agent version to deploy (e.g. \"1.58.0\").",
						},
						// Single config version per agent. Sent to the API as a
						// one-element list; the API type is a slice but we
						// intentionally expose only one version per agent block.
						"configuration_version_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Configuration version entity GUID to associate with this agent in the deployment.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags for the deployment in format 'key:value1,value2'.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. Auto-fetched from the account if not provided.",
			},
			// Computed
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity GUID of the deployment.",
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current phase of the deployment (e.g. CREATED, IN_PROGRESS, FAILED, COMPLETED).",
			},
		},
	}
}

// resourceNewRelicFleetDeploymentCustomizeDiff runs two plan-time validations:
//  1. No two agent blocks may declare the same agent_type.
//  2. If the deployment already exists and its phase is not CREATED, any
//     attempt to change mutable fields is rejected early with a clear error
//     rather than surfacing an opaque API error at apply time.
func resourceNewRelicFleetDeploymentCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Validate duplicate agent types.
	if agentsRaw, ok := d.GetOk("agent"); ok {
		agents := agentsRaw.([]interface{})
		seen := make(map[string]int, len(agents))
		for i, raw := range agents {
			m := raw.(map[string]interface{})
			agentType := m["agent_type"].(string)
			if prev, exists := seen[agentType]; exists {
				return fmt.Errorf(
					"duplicate agent_type %q: agent blocks at index %d and %d both declare the same type — each agent type may appear at most once per deployment",
					agentType, prev, i,
				)
			}
			seen[agentType] = i
		}
	}

	// Phase-gate: block updates once the deployment has left the CREATED phase.
	// d.Id() is non-empty only when the resource already exists in state.
	if d.Id() != "" && d.HasChanges("name", "description", "agent", "tags") {
		phase := d.Get("phase").(string)
		if phase != "" && phase != "CREATED" {
			return fmt.Errorf(
				"cannot update fleet deployment: the deployment is in phase %q and can only be updated while in phase CREATED. "+
					"To make changes, destroy this deployment and create a new one.",
				phase,
			)
		}
	}

	return nil
}

func resourceNewRelicFleetDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	agents, err := expandFleetDeploymentAgents(d.Get("agent").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := parseFleetTags(d.Get("tags").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	input := fleetcontrol.FleetControlFleetDeploymentCreateInput{
		FleetId:     d.Get("fleet_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Agents:      agents,
		Tags:        tags,
		Scope: fleetcontrol.FleetControlScopedReferenceInput{
			ID:   organizationID,
			Type: fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION,
		},
	}

	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateFleetDeploymentWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating fleet deployment: %w", err))
	}

	d.SetId(result.Entity.ID)
	_ = d.Set("deployment_id", result.Entity.ID)
	_ = d.Set("organization_id", organizationID)

	log.Printf("[DEBUG] Created fleet deployment: %s", result.Entity.ID)

	return resourceNewRelicFleetDeploymentRead(ctx, d, meta)
}

func resourceNewRelicFleetDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	entityInterface, err := providerConfig.NewClient.FleetControl.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading fleet deployment %s: %w", d.Id(), err))
	}

	if entityInterface == nil {
		d.SetId("")
		return nil
	}

	entity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetDeploymentEntity)
	if !ok {
		return diag.Errorf("entity '%s' is not a fleet deployment", d.Id())
	}

	_ = d.Set("deployment_id", entity.ID)
	_ = d.Set("fleet_id", entity.FleetId)
	_ = d.Set("phase", string(entity.Phase))

	// name and description are in the query but the API may return empty strings
	// for deployments that have them set; only overwrite state when the API
	// actually returns a value to prevent wiping user-configured values.
	if entity.Name != "" {
		_ = d.Set("name", entity.Name)
	}
	if entity.Description != "" {
		_ = d.Set("description", entity.Description)
	}

	if entity.Scope.ID != "" {
		_ = d.Set("organization_id", entity.Scope.ID)
	}

	// agents is absent from the GetEntity query fragment for
	// EntityManagementFleetDeploymentEntity — entity.Agents is always nil.
	if len(entity.Agents) > 0 {
		if err := d.Set("agent", flattenFleetDeploymentAgents(entity.Agents)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting agent: %w", err))
		}
	}

	// The deployment API accepts tags in the mutation but does not persist or
	// return them via GetEntity; skip the set when the API returns nothing to
	// avoid wiping user-configured tags from state on every refresh.
	if len(entity.Tags) > 0 {
		if err := d.Set("tags", flattenFleetTags(entity.Tags)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
		}
	}

	return nil
}

func resourceNewRelicFleetDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !d.HasChanges("name", "description", "agent", "tags") {
		return nil
	}

	input := fleetcontrol.FleetControlFleetDeploymentUpdateInput{}

	// Always send name and description together to keep the object consistent on
	// the API side — the top-level HasChanges guard already ensures we only
	// reach here when something actually changed.
	input.Name = d.Get("name").(string)
	input.Description = d.Get("description").(string)

	if d.HasChange("agent") {
		agents, err := expandFleetDeploymentAgents(d.Get("agent").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		input.Agents = agents
	}

	if d.HasChange("tags") {
		tags, err := parseFleetTags(d.Get("tags").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		input.Tags = tags
	}

	_, err := providerConfig.NewClient.FleetControl.FleetControlUpdateFleetDeploymentWithContext(ctx, input, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating fleet deployment %s: %w", d.Id(), err))
	}

	return resourceNewRelicFleetDeploymentRead(ctx, d, meta)
}

func resourceNewRelicFleetDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteFleetDeploymentWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting fleet deployment %s: %w", d.Id(), err))
	}

	log.Printf("[DEBUG] Deleted fleet deployment: %s", d.Id())
	return nil
}

// expandFleetDeploymentAgents converts the agent list from schema into API input structs.
// configuration_version_id is a single string in the schema but sent to the
// API as a one-element list.
func expandFleetDeploymentAgents(raw []interface{}) ([]fleetcontrol.FleetControlAgentInput, error) {
	agents := make([]fleetcontrol.FleetControlAgentInput, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})

		agent := fleetcontrol.FleetControlAgentInput{
			AgentType: m["agent_type"].(string),
			Version:   m["version"].(string),
		}

		if v, ok := m["configuration_version_id"].(string); ok && v != "" {
			agent.ConfigurationVersionList = []fleetcontrol.FleetControlConfigurationVersionListInput{
				{ID: v},
			}
		}

		agents = append(agents, agent)
	}
	return agents, nil
}

// flattenFleetDeploymentAgents converts API agent structs back into schema-compatible maps.
// Only the first configuration version is surfaced (matching the single-version schema).
func flattenFleetDeploymentAgents(agents []fleetcontrol.EntityManagementAgentToDeploy) []interface{} {
	result := make([]interface{}, 0, len(agents))
	for _, a := range agents {
		configVersionID := ""
		if len(a.ConfigurationVersionList) > 0 {
			configVersionID = a.ConfigurationVersionList[0].ID
		}

		result = append(result, map[string]interface{}{
			"agent_type":               a.AgentType,
			"version":                  a.Version,
			"configuration_version_id": configVersionID,
		})
	}
	return result
}

// flattenFleetTags converts API tag structs back into "key:value1,value2" strings.
func flattenFleetTags(tags []fleetcontrol.EntityManagementTag) []string {
	result := make([]string, 0, len(tags))
	for _, t := range tags {
		result = append(result, fmt.Sprintf("%s:%s", t.Key, strings.Join(t.Values, ",")))
	}
	return result
}
