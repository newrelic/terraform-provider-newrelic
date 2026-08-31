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
				Description: "The customer ID associated with the organization.",
			},
			"new_managed_account": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Attributes for creating a new managed account alongside the organization.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name of the new account to be created.",
						},
						"region_code": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(listValidOrganizationRegionCodes(), false),
							Description:  fmt.Sprintf("The region code for the account. One of: (%s).", listValidOrganizationRegionCodesStr()),
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
							Description: "The limiting role ID the new organization will be granted for the shared account.",
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
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the organization.",
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

func listValidOrganizationRegionCodesStr() string {
	codes := listValidOrganizationRegionCodes()
	result := ""
	for i, c := range codes {
		if i > 0 {
			result += ", "
		}
		result += c
	}
	return result
}

var (
	_ = context.Background
	_ = log.Printf
	_ = nrErrors.NotFound{}
	_ = diag.FromErr
)

func resourceNewRelicOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	customerID, newManagedAccount, org, sharedAccount, err := expandOrganization(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic organization %s", org.Name)

	result, err := client.Organization.OrganizationCreateWithContext(ctx, customerID, newManagedAccount, org, sharedAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.JobId)
	_ = d.Set("job_id", result.JobId)

	return nil
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

	organizationID := d.Get("org_id").(string)

	updateInput := organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}

	log.Printf("[INFO] Updating New Relic organization %s", organizationID)

	result, err := client.Organization.OrganizationUpdateWithContext(ctx, updateInput, organizationID)
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
	log.Printf("[INFO] Deleting New Relic organization %s is not supported via API; removing from state", d.Id())
	return nil
}
