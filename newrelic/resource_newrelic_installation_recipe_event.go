package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/installevents"
)

func resourceNewRelicInstallationRecipeEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicInstallationRecipeEventCreate,
		ReadContext:   resourceNewRelicInstallationRecipeEventRead,
		UpdateContext: resourceNewRelicInstallationRecipeEventUpdate,
		DeleteContext: resourceNewRelicInstallationRecipeEventDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cli_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"complete": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"details": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"optimized_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"host_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"install_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_library_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kernel_arch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kernel_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_family": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AVAILABLE",
					"CANCELED",
					"DETECTED",
					"FAILED",
					"INSTALLED",
					"INSTALLING",
					"RECOMMENDED",
					"SKIPPED",
					"UNSUPPORTED",
				}, false),
			},
			"targeted_install": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"task_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"validation_duration_milliseconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

var _ = context.Background
var _ = (*installevents.InstallationRecipeStatus)(nil)
var _ diag.Diagnostics(nil)
func resourceNewRelicInstallationRecipeEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input, err := expandInstallationRecipeEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Installevents.InstallationCreateRecipeEventWithContext(ctx, accountID, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.InstallId)
	return diag.FromErr(flattenInstallationRecipeEvent(result, d))
}

func resourceNewRelicInstallationRecipeEventRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicInstallationRecipeEventUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input, err := expandInstallationRecipeEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Installevents.InstallationCreateRecipeEventWithContext(ctx, accountID, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenInstallationRecipeEvent(result, d))
}

func resourceNewRelicInstallationRecipeEventDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}