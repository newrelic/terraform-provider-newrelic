package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
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
				Optional:    true,
				Description: "The private location description.",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the private location.",
				Required:    true,
			},
			"verifiedScriptExecution": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "The private location requires a password to edit if value is true.",
			},
			"domainId": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private location globally unique identifier.",
			},
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the entity to tag.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " The private locations key.",
			},
			"locationId": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
	verifiedScriptExecution := d.Get("verifiedScriptExecution").(bool)
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

	_ = d.Set("guid", res.GUID)

	return nil
}

func resourceNewRelicSyntheticsPrivateLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//
	return nil
}

func resourceNewRelicSyntheticsPrivateLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics
	description := d.Get("description").(string)
	guid := synthetics.EntityGUID(d.Id())
	verifiedScriptExecution := d.Get("verifiedScriptExecution").(bool)
	res, err := client.Synthetics.SyntheticsUpdatePrivateLocation(description, guid, verifiedScriptExecution)

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

	return nil
}

func resourceNewRelicSyntheticsPrivateLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	var diags diag.Diagnostics
	guid := synthetics.EntityGUID(d.Id())
	res, err := client.Synthetics.SyntheticsDeletePrivateLocationWithContext(ctx, guid)

	if err != nil {
		return diag.FromErr(err) //delete error return
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
