package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

// Schema attribute name constants. Use these everywhere instead of bare string
// literals so that a rename only requires a change in one place.
const (
	nrqlMacroAttrName           = "name"
	nrqlMacroAttrExpression     = "expression"
	nrqlMacroAttrOrganizationID = "organization_id"
	nrqlMacroAttrParameterCount = "parameter_count"
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
			nrqlMacroAttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique name for the macro within your organization. Once created, the name cannot be changed — it is the identifier used to reference the macro as `$name` in NRQL queries.",
			},
			nrqlMacroAttrExpression: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The NRQL fragment this macro expands to wherever `$name` appears in a query. Changes to this field take effect immediately across every dashboard widget, alert condition, and query that references the macro.",
			},
			nrqlMacroAttrOrganizationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The GUID of the organization that owns this macro. Populated automatically from your provider credentials — NRQL macros are always scoped to your organization.",
			},
			nrqlMacroAttrParameterCount: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of positional parameters the expression accepts, expressed as bare `$1`, `$2`, ... tokens. Required when the expression contains any positional parameters; defaults to `0` for expressions that do not.",
			},
		},
	}
}

func resourceNewRelicNrqlMacroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	orgID, err := getOrganizationID(ctx, providerConfig, "")
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(nrqlMacroAttrName).(string)
	expression := d.Get(nrqlMacroAttrExpression).(string)
	parameterCount := d.Get(nrqlMacroAttrParameterCount).(int)

	createInput := nrqlmacros.EntityManagementNRQLMacroEntityCreateInput{
		Name:          name,
		Expression:    expression,
		NumParameters: parameterCount,
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
	_ = d.Set(nrqlMacroAttrOrganizationID, orgID)
	_ = d.Set(nrqlMacroAttrName, created.Entity.Name)
	_ = d.Set(nrqlMacroAttrExpression, created.Entity.Expression)
	_ = d.Set(nrqlMacroAttrParameterCount, created.Entity.NumParameters)

	return nil
}

func resourceNewRelicNrqlMacroRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	macroID := d.Id()
	log.Printf("[INFO] Reading New Relic NRQL Macro %s", macroID)

	resp, err := client.Nrqlmacros.GetEntityWithContext(ctx, macroID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] NRQL macro %s not found; removing from state", macroID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	switch entityType := (*resp).(type) {
	case *nrqlmacros.EntityManagementNRQLMacroEntity:
		entity := (*resp).(*nrqlmacros.EntityManagementNRQLMacroEntity)

		if err := d.Set(nrqlMacroAttrName, entity.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(nrqlMacroAttrExpression, entity.Expression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(nrqlMacroAttrOrganizationID, entity.Scope.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(nrqlMacroAttrParameterCount, entity.NumParameters); err != nil {
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
	updateInput := nrqlmacros.EntityManagementNRQLMacroEntityUpdateInput{}

	if d.HasChange(nrqlMacroAttrExpression) {
		updateInput.Expression = d.Get(nrqlMacroAttrExpression).(string)
	}
	if d.HasChange(nrqlMacroAttrParameterCount) {
		n := d.Get(nrqlMacroAttrParameterCount).(int)
		updateInput.NumParameters = &n
	}

	log.Printf("[INFO] Updating New Relic NRQL Macro %s", macroID)

	updated, err := client.Nrqlmacros.EntityManagementUpdateNRQLMacroWithContext(ctx, macroID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if updated.Entity.ID == "" {
		return diag.FromErr(fmt.Errorf("error updating NRQL macro with ID %s", macroID))
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
		return diag.FromErr(fmt.Errorf("error deleting NRQL macro with ID %s", macroID))
	}

	return nil
}
