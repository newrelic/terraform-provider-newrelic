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
				Optional:    true,
				ForceNew:    true,
				Description: "The entity GUID to associate with the change tracking event via entitySearch.",
			},
			"entity_search_query": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Entity search query string that matches exactly one entity.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "A link to the system that generated the change event.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the event.",
			},
			"short_description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A short description of the change suitable for showing in various parts of New Relic One.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An identifier used to correlate account-wide changes across entities.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name or identifier of the person responsible for the change.",
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The category of the change event.",
			},
			"category_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidChangeTrackingCategoryTypes(), false),
				Description:  fmt.Sprintf("The category and type of the change event. One of: (%s).", "see listValidChangeTrackingCategoryTypes()"),
			},
			"custom_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A custom type override. Required when category_type is CUSTOMER_DEFINED__CUSTOM.",
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidChangeTrackingDeploymentTypes(), false),
				Description:  fmt.Sprintf("The type of deployment. One of: (%s).", "see listValidChangeTrackingDeploymentTypes()"),
			},
			"feature_flag_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The identifier of the feature flag.",
			},
			"validation_flags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Flags for validation.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(listValidChangeTrackingValidationFlags(), false),
				},
			},
		},
	}
}

func listValidChangeTrackingCategoryTypes() []string {
	return []string{
		string(changetracking.ChangeTrackingCategoryTypeTypes.BUSINESS_EVENT__CONVENTION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.BUSINESS_EVENT__MARKETING_CAMPAIGN),
		string(changetracking.ChangeTrackingCategoryTypeTypes.BUSINESS_EVENT__OTHER),
		string(changetracking.ChangeTrackingCategoryTypeTypes.CUSTOMER_DEFINED__CUSTOM),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__ARTIFACT_COPY),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__ARTIFACT_DELETION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__ARTIFACT_DEPLOYMENT),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__ARTIFACT_MOVE),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__BUILD_DELETION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__BUILD_PROMOTION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__BUILD_UPLOAD),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__IMAGE_DELETION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__IMAGE_PROMOTION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__IMAGE_PUSH),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_CREATION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_DELETION),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT_LIFECYCLE__RELEASE_BUNDLE_SIGN),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__BASIC),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__BLUE_GREEN),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__CANARY),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__OTHER),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__ROLLING),
		string(changetracking.ChangeTrackingCategoryTypeTypes.DEPLOYMENT__SHADOW),
		string(changetracking.ChangeTrackingCategoryTypeTypes.FEATURE_FLAG__BASIC),
		string(changetracking.ChangeTrackingCategoryTypeTypes.OPERATIONAL__CRASH),
		string(changetracking.ChangeTrackingCategoryTypeTypes.OPERATIONAL__OTHER),
		string(changetracking.ChangeTrackingCategoryTypeTypes.OPERATIONAL__SCHEDULED_MAINTENANCE_PERIOD),
		string(changetracking.ChangeTrackingCategoryTypeTypes.OPERATIONAL__SERVER_REBOOT),
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

func listValidChangeTrackingValidationFlags() []string {
	return []string{
		string(changetracking.ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE),
		string(changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH),
		string(changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_REST_API_FAILURES),
	}
}

var _ = context.Background
var _ = log.Printf
var _ = diag.FromErr
var _ = common.EntityGUID("")

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
		changeTrackingID = fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int))
	}

	d.SetId(changeTrackingID)
	return nil
}
