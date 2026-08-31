package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description: "The customer ID for the organization.",
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
				Description: "Attributes for sharing an account with the new organization.",
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
							Description: "The limiting role ID the new organization will be granted for the shared account.",
						},
					},
				},
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job ID of the organization creation task.",
			},
		},
	}
}
func resourceNewRelicOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	organizationInput := organization.OrganizationCreateOrganizationInput{
		Name: d.Get("name").(string),
	}

	customerID := d.Get("customer_id").(string)

	var newManagedAccountInput *organization.OrganizationNewManagedAccountInput
	if v, ok := d.GetOk("new_managed_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			newManagedAccountInput = &organization.OrganizationNewManagedAccountInput{}
			if n, ok := cfg["name"].(string); ok && n != "" {
				newManagedAccountInput.Name = n
			}
			if rc, ok := cfg["region_code"].(string); ok && rc != "" {
				newManagedAccountInput.RegionCode = organization.OrganizationRegionCodeEnum(rc)
			}
		}
	}

	var sharedAccountInput *organization.OrganizationSharedAccountInput
	if v, ok := d.GetOk("shared_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			sharedAccountInput = &organization.OrganizationSharedAccountInput{
				AccountID: cfg["account_id"].(int),
			}
			if lri, ok := cfg["limiting_role_id"].(int); ok {
				sharedAccountInput.LimitingRoleId = lri
			}
		}
	}

	result, err := client.OrganizationManagement.OrganizationCreateWithContext(ctx, customerID, newManagedAccountInput, organizationInput, sharedAccountInput)
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

	result, err := client.OrganizationManagement.GetOrganizationWithContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil || result.ID == "" {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", result.Name)

	return nil
}

func resourceNewRelicOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	updateInput := organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}

	result, err := client.OrganizationManagement.OrganizationUpdateWithContext(ctx, updateInput, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(result.Errors) > 0 {
		var diags diag.Diagnostics
		for _, e := range result.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  string(e.Type),
				Detail:   e.Message,
			})
		}
		return diags
	}

	_ = d.Set("name", result.OrganizationInformation.Name)

	return resourceNewRelicOrganizationRead(ctx, d, meta)
}

func resourceNewRelicOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}