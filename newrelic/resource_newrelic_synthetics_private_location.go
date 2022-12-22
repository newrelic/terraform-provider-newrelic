package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func resourceNewRelicSyntheticsPrivateLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsPrivateLocationCreate,
		ReadContext:   resourceNewRelicSyntheticsPrivateLocationRead,
		UpdateContext: resourceNewRelicSyntheticsPrivateLocationUpdate,
		DeleteContext: resourceNewRelicSyntheticsPrivateLocationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the account in New Relic.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private location description.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the private location.",
				ForceNew:    true,
				Required:    true,
			},
			"verified_script_execution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The private location requires a password to edit if value is true.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private location globally unique identifier.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The guid of the entity to tag.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private locations key.",
			},
			"location_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An alternate identifier based on name.",
			},
		},
	}
}

func resourceNewRelicSyntheticsPrivateLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	var diags diag.Diagnostics

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	verifiedScriptExecution := d.Get("verified_script_execution").(bool)
	res, err := client.Synthetics.SyntheticsCreatePrivateLocationWithContext(ctx, accountID, description, name, verifiedScriptExecution)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	d.SetId(string(res.GUID))

	_ = d.Set("domain_id", res.DomainId)
	_ = d.Set("key", res.Key)
	_ = d.Set("location_id", res.LocationId)
	_ = d.Set("guid", string(res.GUID))

	return resourceNewRelicSyntheticsPrivateLocationRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsPrivateLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Reading New Relic Synthetics Private Location %s", d.Id())

	guid := common.EntityGUID(d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, guid)
	if err != nil {
		if err.Error() == "Argument \"guid\" has invalid value $guid." {
			return diag.FromErr(fmt.Errorf("invalid GUID"))
		}
		return diag.FromErr(err)
	}
	if _, ok := err.(*errors.NotFound); ok {
		d.SetId("")
		return nil
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	setCommonSyntheticsPrivateLocationAttributes(resp, d)

	return nil
}

func setCommonSyntheticsPrivateLocationAttributes(v *entities.EntityInterface, d *schema.ResourceData) {
	switch e := (*v).(type) {
	case *entities.GenericEntity:
		_ = d.Set("account_id", e.AccountID)
		_ = d.Set("guid", string(e.GUID))
		_ = d.Set("name", e.Name)
	}
}

func resourceNewRelicSyntheticsPrivateLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics

	description := d.Get("description").(string)
	guid := synthetics.EntityGUID(d.Id())
	verifiedScriptExecution := d.Get("verified_script_execution").(bool)

	res, err := client.Synthetics.SyntheticsUpdatePrivateLocationWithContext(ctx, description, guid, verifiedScriptExecution)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	_ = d.Set("domain_id", res.DomainId)
	_ = d.Set("key", res.Key)
	_ = d.Set("location_id", res.LocationId)
	_ = d.Set("guid", string(res.GUID))

	return resourceNewRelicSyntheticsPrivateLocationRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsPrivateLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics
	guid := synthetics.EntityGUID(d.Id())

	res, err := client.Synthetics.SyntheticsDeletePrivateLocationWithContext(ctx, guid)

	if err != nil {
		return diag.FromErr(err)
	}
	if res != nil {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	d.SetId("")
	return nil
}
