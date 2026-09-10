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

func resourceNewRelicChangeTrackingEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicChangeTrackingEventCreate,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"entity_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The entity GUID to associate the change tracking event with.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The version of the deployed software.",
			},
			"changelog": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A URL to the changelog or a list of changes.",
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
				Description: "A link to the system that generated the deployment.",
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidChangeTrackingDeploymentTypes(), false),
				Description:  fmt.Sprintf("The type of deployment. One of: (%s).", listValidChangeTrackingDeploymentTypesString()),
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
				Description: "The start time of the deployment as the number of milliseconds since the Unix epoch.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The username of the deployer or bot.",
			},
		},
	}
}

func listValidChangeTrackingDeploymentTypes() []string {
	return []string{
		string(changetracking.ChangeTrackingDeploymentTypeTypes.BASIC),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.BLUE_GREEN),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.CANARY),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.OTHER),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.ROLLING),
		string(changetracking.ChangeTrackingDeploymentTypeTypes.SHADOW),
	}
}

func listValidChangeTrackingDeploymentTypesString() string {
	types := listValidChangeTrackingDeploymentTypes()
	result := ""
	for i, t := range types {
		if i > 0 {
			result += ", "
		}
		result += t
	}
	return result
}

var _ = log.Printf
var _ = context.Background
var _ *common.EntityGUID
var _ *changetracking.ChangeTrackingDeploymentInput
var _ diag.Diagnostics

func resourceNewRelicChangeTrackingEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	changeTrackingEvent, dataHandlingRules, err := expandChangeTrackingEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic change tracking event")

	result, err := client.ChangeTracking.ChangeTrackingCreateEventWithContext(ctx, changeTrackingEvent, dataHandlingRules)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil {
		return diag.Errorf("error creating change tracking event: empty response")
	}

	if result.ChangeTrackingEvent != nil {
		switch e := result.ChangeTrackingEvent.(type) {
		case *changetracking.ChangeTrackingDeploymentEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingFeatureFlagEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingGenericEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingEvent:
			d.SetId(e.ChangeTrackingId)
		default:
			d.SetId(fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int)))
		}
	} else {
		d.SetId(fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int)))
	}

	return nil
}
