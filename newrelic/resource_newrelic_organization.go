package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
)

func resourceNewRelicOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicOrganizationCreate,
		ReadContext:   resourceNewRelicOrganizationRead,
		UpdateContext: resourceNewRelicOrganizationUpdate,
		DeleteContext: resourceNewRelicOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The customer ID for the organization.",
			},
			"new_managed_account": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Attributes for a new managed account to create alongside the organization.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name of the new managed account to be created.",
						},
						"region_code": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(listValidOrganizationRegionCodes(), false),
							Description:  fmt.Sprintf("The region code for the account to be created. One of: (%s).", joinStrings(listValidOrganizationRegionCodes())),
						},
					},
				},
			},
			"shared_account": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Attributes for sharing an account with the new organization.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of the account to share with the new organization.",
						},
						"limiting_role_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The limiting role ID the new organization will be granted on the shared account.",
						},
					},
				},
			},
			// Computed
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job ID of the organization creation task.",
			},
		},
	}
}

func listValidOrganizationRegionCodes() []string {
	return []string{
		string(organization.OrganizationRegionCodeEnumTypes.EU01),
		string(organization.OrganizationRegionCodeEnumTypes.US01),
	}
}

func joinStrings(s []string) string {
	result := ""
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}

var _ = context.Background
var _ = log.Printf
var _ = nrErrors.NotFound{}

func resourceNewRelicOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	input, err := expandOrganization(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic organization %s", input.Organization.Name)

	result, err := client.Organization.OrganizationCreateWithContext(
		ctx,
		input.CustomerID,
		input.NewManagedAccount,
		input.Organization,
		input.SharedAccount,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.JobId)
	_ = d.Set("job_id", result.JobId)

	return resourceNewRelicOrganizationRead(ctx, d, meta)
}

func resourceNewRelicOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic organization %s", d.Id())

	result, err := client.Organization.GetOrganizationWithContext(ctx)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenOrganization(result, d))
}

func resourceNewRelicOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Updating New Relic organization %s", d.Id())

	updateInput := organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}

	result, err := client.Organization.OrganizationUpdateWithContext(ctx, updateInput, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(result.Errors) > 0 {
		var diagErrors diag.Diagnostics
		for _, e := range result.Errors {
			diagErrors = append(diagErrors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  string(e.Type),
				Detail:   e.Message,
			})
		}
		return diagErrors
	}

	return resourceNewRelicOrganizationRead(ctx, d, meta)
}

func resourceNewRelicOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting New Relic organization is not supported via API; removing from state only.")
	return nil
}
