package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

func resourceNewRelicFleet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetCreate,
		ReadContext:   resourceNewRelicFleetRead,
		UpdateContext: resourceNewRelicFleetUpdate,
		DeleteContext: resourceNewRelicFleetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the fleet.",
			},
			"managed_entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(fleetcontrol.EntityManagementManagedEntityTypeTypes.HOST),
					string(fleetcontrol.EntityManagementManagedEntityTypeTypes.KUBERNETESCLUSTER),
				}, false),
				Description: "The type of entities this fleet will manage. Allowed values: HOST, KUBERNETESCLUSTER.",
			},
			"operating_system": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LINUX",
					"WINDOWS",
				}, false),
				Description: "The operating system type. Required for HOST fleets. Allowed values: LINUX, WINDOWS. Must not be set for KUBERNETESCLUSTER fleets.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the fleet.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags for the fleet in format 'key:value1,value2'. Each tag can have multiple values.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization ID. If not provided, it will be auto-fetched from the account.",
			},
		},
	}
}

func resourceNewRelicFleetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get or fetch organization ID
	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate operating system requirements
	managedEntityType := d.Get("managed_entity_type").(string)
	operatingSystem, hasOS := d.GetOk("operating_system")

	if managedEntityType == "HOST" && !hasOS {
		return diag.Errorf("operating_system is required when managed_entity_type is HOST")
	}
	if managedEntityType == "KUBERNETESCLUSTER" && hasOS {
		return diag.Errorf("operating_system should not be specified for KUBERNETESCLUSTER fleets")
	}

	// Map managed entity type
	entityType, err := mapManagedEntityType(managedEntityType)
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse tags
	tags, err := parseFleetTags(d.Get("tags").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Build create input
	createInput := fleetcontrol.FleetControlFleetEntityCreateInput{
		Name:              d.Get("name").(string),
		ManagedEntityType: entityType,
		Scope: fleetcontrol.FleetControlScopedReferenceInput{
			ID:   organizationID,
			Type: fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION,
		},
	}

	// Add optional fields
	if desc, ok := d.GetOk("description"); ok {
		createInput.Description = desc.(string)
	}

	// Add operating system if provided
	if hasOS {
		osType, osErr := mapOperatingSystemType(operatingSystem.(string))
		if osErr != nil {
			return diag.FromErr(osErr)
		}
		createInput.OperatingSystem = &fleetcontrol.FleetControlOperatingSystemCreateInput{
			Type: osType,
		}
	}

	if len(tags) > 0 {
		createInput.Tags = tags
	}

	// Create fleet
	result, err := providerConfig.NewClient.FleetControl.FleetControlCreateFleetWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Entity.ID)

	// Flatten the result
	if err := flattenFleetControlEntity(&result.Entity, d, organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Get organization ID from state, or fetch if not present (e.g., during import)
	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		var err error
		organizationID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Fetch fleet entity
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

	// Type assert to fleet entity
	fleetEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetEntity)
	if !ok {
		return diag.Errorf("entity '%s' is not a fleet", d.Id())
	}

	// Flatten into state
	if err := flattenFleetEntity(fleetEntity, d, organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	// Check if any updateable fields changed
	if !d.HasChanges("name", "description", "tags") {
		return nil
	}

	// Get organization ID
	organizationID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Build update input
	updateInput := fleetcontrol.FleetControlUpdateFleetEntityInput{}

	if d.HasChange("name") {
		updateInput.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateInput.Description = d.Get("description").(string)
	}

	if d.HasChange("tags") {
		tags, tagsErr := parseFleetTags(d.Get("tags").([]interface{}))
		if tagsErr != nil {
			return diag.FromErr(tagsErr)
		}
		updateInput.Tags = tags
	}

	// Update fleet
	result, err := providerConfig.NewClient.FleetControl.FleetControlUpdateFleetWithContext(ctx, updateInput, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Flatten the result
	if err := flattenFleetControlEntity(&result.Entity, d, organizationID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	log.Printf("[INFO] Deleting New Relic Fleet %s", d.Id())

	_, err := providerConfig.NewClient.FleetControl.FleetControlDeleteFleetWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
