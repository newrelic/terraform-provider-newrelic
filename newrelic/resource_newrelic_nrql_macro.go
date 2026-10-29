package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

func resourceNewRelicNrqlMacro() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNrqlMacroCreate,
		ReadContext:   resourceNewRelicNrqlMacroRead,
		UpdateContext: resourceNewRelicNrqlMacroUpdate,
		DeleteContext: resourceNewRelicNrqlMacroDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique name of the NRQL macro. Used to reference it as $name in queries. Immutable after creation.",
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NRQL expression fragment that this macro expands to when referenced in a query. Updating this value will immediately affect all queries, dashboards, and alert conditions that reference this macro.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The organization GUID to which this macro belongs. If omitted, the organization ID is resolved from the provider credentials.",
			},
			"num_parameters": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of positional parameters ($1, $2, ...) used in the expression.",
			},
		},
	}
}

func resourceNewRelicNrqlMacroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	expression := d.Get("expression").(string)
	numParameters := d.Get("num_parameters").(int)

	createInput := nrqlmacros.EntityManagementNRQLMacroEntityCreateInput{
		Name:          name,
		Expression:    expression,
		NumParameters: numParameters,
		Scope: nrqlmacros.EntityManagementScopedReferenceInput{
			Type: nrqlmacros.EntityManagementEntityScopeTypes.ORGANIZATION,
			ID:   orgID,
		},
	}

	log.Printf("[INFO] Creating New Relic NRQL Macro: %s", name)

	created, err := client.Nrqlmacros.EntityManagementCreateNRQLMacroWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created.Entity.ID == "" {
		return diag.Errorf("the New Relic API did not return an entity ID for the new NRQL macro")
	}

	log.Printf("[INFO] New Relic NRQL Macro created, ID: %s", created.Entity.ID)

	d.SetId(created.Entity.ID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("name", created.Entity.Name)
	_ = d.Set("expression", created.Entity.Expression)
	_ = d.Set("num_parameters", created.Entity.NumParameters)

	return nil
}

func resourceNewRelicNrqlMacroRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	macroID := d.Id()
	log.Printf("[INFO] Reading New Relic NRQL Macro %s", macroID)

	resp, err := client.Nrqlmacros.GetEntityWithContext(ctx, macroID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		// Entity not found; remove from state so Terraform plans a re-create.
		d.SetId("")
		return nil
	}

	switch entityType := (*resp).(type) {
	case *nrqlmacros.EntityManagementNRQLMacroEntity:
		entity := (*resp).(*nrqlmacros.EntityManagementNRQLMacroEntity)

		if err := d.Set("name", entity.Name); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("expression", entity.Expression); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("organization_id", entity.Scope.ID); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("num_parameters", entity.NumParameters); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}

	return nil
}

func resourceNewRelicNrqlMacroUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	macroID := d.Id()

	updateInput := nrqlmacros.EntityManagementNRQLMacroEntityUpdateInput{
		Expression:    d.Get("expression").(string),
		NumParameters: d.Get("num_parameters").(int),
	}

	log.Printf("[INFO] Updating New Relic NRQL Macro %s", macroID)

	updated, err := client.Nrqlmacros.EntityManagementUpdateNRQLMacroWithContext(ctx, macroID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if updated.Entity.ID == "" {
		return diag.FromErr(fmt.Errorf("error in updating NRQL macro with ID %s", macroID))
	}

	return resourceNewRelicNrqlMacroRead(ctx, d, meta)
}

func resourceNewRelicNrqlMacroDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	macroID := d.Id()
	log.Printf("[INFO] Deleting New Relic NRQL Macro %s", macroID)

	result, err := client.Nrqlmacros.EntityManagementDelete(macroID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.ID == "" {
		return diag.FromErr(fmt.Errorf("error in deleting NRQL macro with ID %s", macroID))
	}

	return nil
}
