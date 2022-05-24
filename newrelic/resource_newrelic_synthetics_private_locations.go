package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
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
	pl := expandSyntheticsPrivateLocation(d)

	var diags diag.Diagnostics

	res, err := client.Synthetics.SyntheticsCreatePrivateLocationWithContext(ctx,accountID,"description","syn",true)
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

	d.SetId(res.Key)

	return resourceNewRelicSyntheticsPrivateLocationRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsPrivateLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = %s", d.Id())

	//entityResults, err := client.Entities.GetEntitySearchWithContext(ctx, entities.EntitySearchOptions{}, queryString, entities.EntitySearchQueryBuilder{}, []entities.EntitySearchSortCriteria{})
	//if err != nil {
	//
	//}
	//var entity *entities.EntityOutlineInterface
	//for _, e := range entityResults.Results.Entities {
	//	if e.GetName() == d.Id() {
	//		entity = &e
	//		break
	//	}
	//}

	return nil
}

func expandSyntheticsPrivateLocation(d *schema.ResourceData) *synthetics.SecureCredential {

	pl := synthetics.
		Description:             d.Get("description").(string),
		Name:                    d.Get(),
		VerifiedScriptExecution: d.Get("boolean").(true),
	}

	return &pl
}

func flattenSyntheticsPrivateLocation(pl *synthetics.SyntheticsPrivateLocationMutationResult, d *schema.ResourceData) error {

	return nil
}
