package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the fleet to deploy to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the deployment.",
			},
			"agent": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Agent configurations for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of agent (e.g., NRInfra, NRDOT, FluentBit, NRPrometheusAgent).",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The agent version (e.g., 1.70.0, 2.0.0, or '*' for KUBERNETESCLUSTER fleets only).",
						},
						"configuration_version_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of configuration version IDs to deploy with this agent.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the deployment.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags for the deployment in format 'key:value1,value2'.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current phase of the deployment.",
			},
		},
	}
}

func resourceNewRelicFleetDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	fleetID := d.Get("fleet_id").(string)

	// Parse agent specifications
	agents, err := parseAgentSpecs(d.Get("agent").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate agent versions for fleet type
	if validationErr := validateAgentVersionsForFleet(ctx, providerConfig, fleetID, agents); validationErr != nil {
		return diag.FromErr(validationErr)
	}

	// Parse tags
	tags, err := parseFleetTags(d.Get("tags").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Build create input
	createInput := fleetcontrol.FleetControlFleetDeploymentCreateInput{
		FleetId: fleetID,
		Name:    d.Get("name").(string),
		Agents:  agents,
	}

	if desc, ok := d.GetOk("description"); ok {
		createInput.Description = desc.(string)
	}

	if len(tags) > 0 {
		createInput.Tags = tags
	}

	// Create deployment
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateFleetDeploymentWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Entity.ID)

	// Flatten the result
	if err := flattenFleetControlDeployment(&result.Entity, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Fetch deployment entity
	entityInterface, err := providerConfig.NewClient.FleetControl.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if entityInterface == nil {
		d.SetId("")
		return nil
	}

	// Type assert to deployment entity
	deploymentEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetDeploymentEntity)
	if !ok {
		return diag.Errorf("entity '%s' is not a fleet deployment", d.Id())
	}

	// Set fleet_id
	if err := d.Set("fleet_id", deploymentEntity.FleetId); err != nil {
		return diag.FromErr(err)
	}

	// Note: For full flattening, we'd need to convert EntityManagement types to FleetControl types
	// For now, set basic fields
	if err := d.Set("name", deploymentEntity.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", deploymentEntity.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("phase", string(deploymentEntity.Phase)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Check if any updateable fields changed
	if !d.HasChanges("name", "description", "agent", "tags") {
		return nil
	}

	// Build update input
	updateInput := fleetcontrol.FleetControlFleetDeploymentUpdateInput{}

	if d.HasChange("name") {
		updateInput.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateInput.Description = d.Get("description").(string)
	}

	if d.HasChange("agent") {
		agents, err := parseAgentSpecs(d.Get("agent").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		// Validate agent versions
		fleetID := d.Get("fleet_id").(string)
		if err := validateAgentVersionsForFleet(ctx, providerConfig, fleetID, agents); err != nil {
			return diag.FromErr(err)
		}
		updateInput.Agents = agents
	}

	if d.HasChange("tags") {
		tags, err := parseFleetTags(d.Get("tags").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateInput.Tags = tags
	}

	// Update deployment
	result, err := providerConfig.NewClient.FleetControl.FleetControlUpdateFleetDeploymentWithContext(ctx, updateInput, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Flatten the result
	if err := flattenFleetControlDeployment(&result.Entity, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	log.Printf("[INFO] Deleting New Relic Fleet Deployment %s", d.Id())

	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteFleetDeploymentWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
