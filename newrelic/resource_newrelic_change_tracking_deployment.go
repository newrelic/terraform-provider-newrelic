package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicChangeTrackingDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicChangeTrackingDeploymentCreate,
		ReadContext:   resourceNewRelicChangeTrackingDeploymentRead,
		DeleteContext: resourceNewRelicChangeTrackingDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_guid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"changelog": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"commit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deep_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BASIC",
					"BLUE_GREEN",
					"CANARY",
					"OTHER",
					"ROLLING",
					"SHADOW",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
func resourceNewRelicChangeTrackingDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	input := changetracking.ChangeTrackingDeploymentInput{
		EntityGUID: common.EntityGUID(d.Get("entity_guid").(string)),
		Version:    d.Get("version").(string),
	}

	if v, ok := d.GetOk("changelog"); ok {
		input.Changelog = v.(string)
	}
	if v, ok := d.GetOk("commit"); ok {
		input.Commit = v.(string)
	}
	if v, ok := d.GetOk("deep_link"); ok {
		input.DeepLink = v.(string)
	}
	if v, ok := d.GetOk("deployment_type"); ok {
		input.DeploymentType = changetracking.ChangeTrackingDeploymentType(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		input.GroupId = v.(string)
	}
	if v, ok := d.GetOk("timestamp"); ok {
		input.Timestamp = nrtime.EpochMilliseconds(v.(int))
	}
	if v, ok := d.GetOk("user"); ok {
		input.User = v.(string)
	}

	dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{}

	result, err := client.ChangeTracking.ChangeTrackingCreateDeploymentWithContext(ctx, dataHandlingRules, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.DeploymentId)
	_ = d.Set("deployment_id", result.DeploymentId)
	_ = d.Set("entity_guid", string(result.EntityGUID))
	_ = d.Set("version", result.Version)
	_ = d.Set("changelog", result.Changelog)
	_ = d.Set("commit", result.Commit)
	_ = d.Set("deep_link", result.DeepLink)
	_ = d.Set("deployment_type", string(result.DeploymentType))
	_ = d.Set("description", result.Description)
	_ = d.Set("group_id", result.GroupId)
	_ = d.Set("timestamp", int(result.Timestamp))
	_ = d.Set("user", result.User)

	return nil
}

func resourceNewRelicChangeTrackingDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicChangeTrackingDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}