package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the organization.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The customer ID to associate with the organization.",
			},
			"new_managed_account": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Attributes for creating a new managed account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the new account to be created.",
						},
						"region_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The region code for the account to be created.",
							ValidateFunc: validation.StringInSlice([]string{
								"EU01",
								"US01",
							}, false),
						},
					},
				},
			},
			"shared_account": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Attributes for creating an account share with the new organization.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The ID of the account to share with the new organization.",
						},
						"limiting_role_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The limiting role ID the new organization will be granted on the shared account.",
						},
					},
				},
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job ID of the organization creation task.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the organization.",
			},
		},
	}
}

var _ = context.Background
var _ = fmt.Sprintf
var _ = log.Printf
var _ = organization.OrganizationRegionCodeEnumTypes
var _ = validation.StringIsNotWhiteSpace
func resourceNewRelicOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	input, err := expandOrganizationCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	customerID := d.Get("customer_id").(string)

	var newManagedAccount *organization.OrganizationNewManagedAccountInput
	if v, ok := d.GetOk("new_managed_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			newManagedAccount = expandOrganizationNewManagedAccount(items[0].(map[string]interface{}))
		}
	}

	var sharedAccount *organization.OrganizationSharedAccountInput
	if v, ok := d.GetOk("shared_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			sharedAccount = expandOrganizationSharedAccount(items[0].(map[string]interface{}))
		}
	}

	result, err := client.Organization.OrganizationCreateWithContext(ctx, customerID, newManagedAccount, *input, sharedAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := flattenOrganizationCreate(result, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicOrganizationRead(ctx, d, meta)
}

func resourceNewRelicOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading NewRelic Organization %s", d.Id())

	result, err := client.Organization.GetOrganizationWithContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil || result.ID == "" {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenOrganization(result, d))
}

func resourceNewRelicOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Updating NewRelic Organization %s", d.Id())

	input, err := expandOrganizationUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Organization.OrganizationUpdateWithContext(ctx, *input, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(result.Errors) > 0 {
		return diag.Errorf("error updating organization: %s", result.Errors[0].Message)
	}

	return resourceNewRelicOrganizationRead(ctx, d, meta)
}

func resourceNewRelicOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting NewRelic Organization is not supported via API; removing from state only.")
	d.SetId("")
	return nil
}