package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

// Ensure unused imports are referenced somewhere in the file (used by CRUD section).
var _ = log.Printf
var _ = context.Background
var _ diag.Diagnostics
var _ common.EntityGUID
var _ nrtime.EpochMilliseconds

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
				Description: "A URL to the changelog or, if not linkable, a list of changes.",
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
				ValidateFunc: validation.StringInSlice(listValidChangeTrackingDeploymentTypes(), false),
				Description:  "The type of deployment.",
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
				ForceNew:    true,
				Description: "The start time of the deployment as the number of milliseconds since the Unix epoch. Defaults to now.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The username of the deployer or bot.",
			},
			// Computed
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique deployment identifier.",
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

func resourceNewRelicChangeTrackingEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Creating New Relic Change Tracking Event")

	changeTrackingEvent, dataHandlingRules, err := expandChangeTrackingEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.ChangeTracking.ChangeTrackingCreateEventWithContext(ctx, changeTrackingEvent, dataHandlingRules)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil {
		return diag.Errorf("error creating change tracking event: empty response")
	}

	var changeTrackingID string
	if result.ChangeTrackingEvent != nil {
		switch e := result.ChangeTrackingEvent.(type) {
		case *changetracking.ChangeTrackingDeploymentEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingFeatureFlagEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingGenericEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingEvent:
			changeTrackingID = e.ChangeTrackingId
		}
	}

	if changeTrackingID == "" {
		changeTrackingID = "change-tracking-event"
	}

	d.SetId(changeTrackingID)
	return nil
}
