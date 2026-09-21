package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/streamingexport"
)

func resourceNewRelicStreamingExportRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicStreamingExportRuleCreate,
		ReadContext:   resourceNewRelicStreamingExportRuleRead,
		DeleteContext: resourceNewRelicStreamingExportRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The account ID of the account that owns the streaming export rule.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the streaming export rule.",
			},
			"nrql": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "NRQL to select the telemetry data to export.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Additional information about the streaming export rule.",
			},
			"payload_compression": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidStreamingExportPayloadCompressionTypes(), false),
				Description:  fmt.Sprintf("Whether to compress payloads before sending them out. One of: (%s).", fmt.Sprintf("%s, %s", string(streamingexport.StreamingExportPayloadCompressionTypes.DISABLED), string(streamingexport.StreamingExportPayloadCompressionTypes.GZIP))),
			},
			"aws": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "AWS parameters for the streaming export rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_account_id": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The AWS account to which the target firehose belongs.",
						},
						"delivery_stream_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the delivery stream to write events to.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The AWS region the delivery stream is located in.",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The role configured for New Relic to assume.",
						},
					},
				},
			},
			"azure": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Azure parameters for the streaming export rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_hub_connection_string": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "Connection string that has access to the specific Event Hub.",
						},
						"event_hub_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of Event Hub to write events to.",
						},
					},
				},
			},
			"gcp": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "GCP parameters for the streaming export rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcp_project_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The GCP project ID.",
						},
						"pubsub_topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The Pub/Sub topic ID.",
						},
					},
				},
			},
			// Computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the streaming export rule.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A message returned by the latest API call.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which the process of creating the streaming rule began.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the status of the streaming rule was updated.",
			},
		},
	}
}

func listValidStreamingExportPayloadCompressionTypes() []string {
	return []string{
		string(streamingexport.StreamingExportPayloadCompressionTypes.DISABLED),
		string(streamingexport.StreamingExportPayloadCompressionTypes.GZIP),
	}
}

// Ensure unused imports are referenced in CRUD functions generated separately.
var _ = context.Background
var _ = log.Printf
var _ = diag.FromErr
var _ = nrErrors.NewNotFound

func resourceNewRelicStreamingExportRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	awsParameters, azureParameters, gcpParameters, ruleParameters, err := expandStreamingExportRule(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic streaming export rule %s", ruleParameters.Name)

	result, err := client.Streamingexport.StreamingExportCreateRuleWithContext(ctx, accountID, awsParameters, azureParameters, gcpParameters, ruleParameters)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))

	return resourceNewRelicStreamingExportRuleRead(ctx, d, meta)
}

func resourceNewRelicStreamingExportRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic streaming export rule %s", d.Id())

	rules, err := client.Streamingexport.GetStreamingRulesWithContext(ctx, accountID)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if rules == nil || len(*rules) == 0 {
		d.SetId("")
		return nil
	}

	var found *streamingexport.StreamingExportRule
	for i := range *rules {
		rule := (*rules)[i]
		if fmt.Sprintf("%d", rule.ID) == d.Id() {
			found = &rule
			break
		}
	}

	if found == nil {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenStreamingExportRule(found, d))
}

func resourceNewRelicStreamingExportRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting New Relic streaming export rule %s", d.Id())
	d.SetId("")
	return nil
}
