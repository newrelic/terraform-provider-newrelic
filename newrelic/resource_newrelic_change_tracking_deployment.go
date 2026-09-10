package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

func resourceNewRelicChangeTrackingDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicChangeTrackingDeploymentCreate,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The version of the deployed software, for example, something like v1.1.",
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The NR entity that was deployed.",
			},
			"changelog": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A URL for the changelog or, if not linkable, a list of changes.",
			},
			"commit": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The commit identifier, for example, a Git commit SHA.",
			},
			"deep_link": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A URL to the system that generated the deployment.",
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidChangeTrackingDeploymentTypeValues(), false),
				Description:  fmt.Sprintf("The type of deployment, for example, 'Blue green' or 'Rolling'. One of: (%s).", listValidChangeTrackingDeploymentTypeValuesString()),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the deployment.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An identifier used to correlate account-wide changes across entities.",
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The username of the deployer or bot.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique deployment identifier.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func listValidChangeTrackingDeploymentTypeValues() []string {
	return []string{
		string(changetracking.ChangeTrackingDeploymentTypeTypes.BASIC),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.BLUE_GREEN),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.CANARY),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.OTHER),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.ROLLING),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.SHADOW),
	}
}

func listValidChangeTrackingDeploymentTypeValuesString() string {
	vals := listValidChangeTrackingDeploymentTypeValues()
	result := ""
	for i, v := range vals {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}

var _ = context.Background
var _ = log.Printf
var _ = common.EntityGUID("")
var _ diag.Diagnostics

func resourceNewRelicChangeTrackingDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	dataHandlingRules, deployment, err := expandChangeTrackingDeployment(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic change tracking deployment for entity %s", deployment.EntityGUID)

	result, err := client.ChangeTracking.ChangeTrackingCreateDeploymentWithContext(ctx, dataHandlingRules, deployment)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.DeploymentId)
	_ = d.Set("deployment_id", result.DeploymentId)

	return nil
}
