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
			// agent blocks define which agent versions (and optionally which
			// configuration versions) are included in this deployment.
			"agent": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "One or more agent blocks defining agent type, version, and optional configuration versions to include in the deployment.",
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
						"configuration_version_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of configuration version entity GUIDs to associate with this agent in the deployment.",
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
				Description: "The current phase of the deployment (e.g. DRAFT, READY).",
			},
		},
	}
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
	_ = d.Set("name", entity.Name)
	_ = d.Set("description", entity.Description)
	_ = d.Set("phase", string(entity.Phase))

	if entity.Scope.ID != "" {
		_ = d.Set("organization_id", entity.Scope.ID)
	}

	if err := d.Set("agent", flattenFleetDeploymentAgents(entity.Agents)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting agent: %w", err))
	}

	if err := d.Set("tags", flattenFleetTags(entity.Tags)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
	}

	return nil
}

func resourceNewRelicFleetDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !d.HasChanges("name", "description", "agent", "tags") {
		return nil
	}

	agents, err := expandFleetDeploymentAgents(d.Get("agent").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := parseFleetTags(d.Get("tags").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	input := fleetcontrol.FleetControlFleetDeploymentUpdateInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Agents:      agents,
		Tags:        tags,
	}

	_, err = providerConfig.NewClient.FleetControl.FleetControlUpdateFleetDeploymentWithContext(ctx, input, d.Id())
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
func expandFleetDeploymentAgents(raw []interface{}) ([]fleetcontrol.FleetControlAgentInput, error) {
	agents := make([]fleetcontrol.FleetControlAgentInput, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})

		agent := fleetcontrol.FleetControlAgentInput{
			AgentType: m["agent_type"].(string),
			Version:   m["version"].(string),
		}

		if v, ok := m["configuration_version_ids"]; ok {
			for _, id := range v.(*schema.Set).List() {
				agent.ConfigurationVersionList = append(
					agent.ConfigurationVersionList,
					fleetcontrol.FleetControlConfigurationVersionListInput{ID: id.(string)},
				)
			}
		}

		agents = append(agents, agent)
	}
	return agents, nil
}

// flattenFleetDeploymentAgents converts API agent structs back into schema-compatible maps.
func flattenFleetDeploymentAgents(agents []fleetcontrol.EntityManagementAgentToDeploy) []interface{} {
	result := make([]interface{}, 0, len(agents))
	for _, a := range agents {
		ids := make([]interface{}, 0, len(a.ConfigurationVersionList))
		for _, cv := range a.ConfigurationVersionList {
			ids = append(ids, cv.ID)
		}

		m := map[string]interface{}{
			"agent_type":                a.AgentType,
			"version":                   a.Version,
			"configuration_version_ids": schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), ids),
		}
		result = append(result, m)
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
