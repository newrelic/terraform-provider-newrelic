package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/streamingexport"
)

func expandStreamingExportRule(d *schema.ResourceData) (streamingexport.StreamingExportAwsInput, streamingexport.StreamingExportAzureInput, streamingexport.StreamingExportGcpInput, streamingexport.StreamingExportRuleInput, error) {
	ruleParameters := streamingexport.StreamingExportRuleInput{
		Name: d.Get("name").(string),
		NRQL: streamingexport.NRQL(d.Get("nrql").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		ruleParameters.Description = v.(string)
	}

	if v, ok := d.GetOk("payload_compression"); ok {
		ruleParameters.PayloadCompression = streamingexport.StreamingExportPayloadCompression(v.(string))
	}

	awsParameters := expandStreamingExportRuleAws(d)
	azureParameters := expandStreamingExportRuleAzure(d)
	gcpParameters := expandStreamingExportRuleGcp(d)

	return awsParameters, azureParameters, gcpParameters, ruleParameters, nil
}

func expandStreamingExportRuleAws(d *schema.ResourceData) streamingexport.StreamingExportAwsInput {
	awsInput := streamingexport.StreamingExportAwsInput{}

	if v, ok := d.GetOk("aws"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			awsInput = expandStreamingExportRuleAwsBlock(cfg)
		}
	}

	return awsInput
}

func expandStreamingExportRuleAwsBlock(cfg map[string]interface{}) streamingexport.StreamingExportAwsInput {
	awsInput := streamingexport.StreamingExportAwsInput{}

	if v, ok := cfg["aws_account_id"]; ok {
		awsInput.AwsAccountId = v.(int)
	}

	if v, ok := cfg["delivery_stream_name"]; ok {
		awsInput.DeliveryStreamName = v.(string)
	}

	if v, ok := cfg["region"]; ok {
		awsInput.Region = v.(string)
	}

	if v, ok := cfg["role"]; ok {
		awsInput.Role = v.(string)
	}

	return awsInput
}

func expandStreamingExportRuleAzure(d *schema.ResourceData) streamingexport.StreamingExportAzureInput {
	azureInput := streamingexport.StreamingExportAzureInput{}

	if v, ok := d.GetOk("azure"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			azureInput = expandStreamingExportRuleAzureBlock(cfg)
		}
	}

	return azureInput
}

func expandStreamingExportRuleAzureBlock(cfg map[string]interface{}) streamingexport.StreamingExportAzureInput {
	azureInput := streamingexport.StreamingExportAzureInput{}

	if v, ok := cfg["event_hub_connection_string"]; ok {
		azureInput.EventHubConnectionString = v.(string)
	}

	if v, ok := cfg["event_hub_name"]; ok {
		azureInput.EventHubName = v.(string)
	}

	return azureInput
}

func expandStreamingExportRuleGcp(d *schema.ResourceData) streamingexport.StreamingExportGcpInput {
	gcpInput := streamingexport.StreamingExportGcpInput{}

	if v, ok := d.GetOk("gcp"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			gcpInput = expandStreamingExportRuleGcpBlock(cfg)
		}
	}

	return gcpInput
}

func expandStreamingExportRuleGcpBlock(cfg map[string]interface{}) streamingexport.StreamingExportGcpInput {
	gcpInput := streamingexport.StreamingExportGcpInput{}

	if v, ok := cfg["gcp_project_id"]; ok {
		gcpInput.GcpProjectId = v.(string)
	}

	if v, ok := cfg["pubsub_topic_id"]; ok {
		gcpInput.PubsubTopicId = v.(string)
	}

	return gcpInput
}

func flattenStreamingExportRule(result *streamingexport.StreamingExportRule, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	if err := d.Set("name", result.Name); err != nil {
		return err
	}

	if err := d.Set("nrql", string(result.NRQL)); err != nil {
		return err
	}

	if err := d.Set("description", result.Description); err != nil {
		return err
	}

	if err := d.Set("payload_compression", string(result.PayloadCompression)); err != nil {
		return err
	}

	if err := d.Set("status", string(result.Status)); err != nil {
		return err
	}

	if err := d.Set("message", result.Message); err != nil {
		return err
	}

	if err := d.Set("created_at", string(result.CreatedAt)); err != nil {
		return err
	}

	if err := d.Set("updated_at", string(result.UpdatedAt)); err != nil {
		return err
	}

	if result.Account.ID != 0 {
		if err := d.Set("account_id", result.Account.ID); err != nil {
			return err
		}
	}

	if result.Aws != (streamingexport.StreamingExportAwsDetails{}) {
		if err := d.Set("aws", flattenStreamingExportAwsDetails(result.Aws)); err != nil {
			return err
		}
	}

	if result.Azure != (streamingexport.StreamingExportAzureDetails{}) {
		if err := d.Set("azure", flattenStreamingExportAzureDetails(result.Azure)); err != nil {
			return err
		}
	}

	if result.Gcp != (streamingexport.StreamingExportGcpDetails{}) {
		if err := d.Set("gcp", flattenStreamingExportGcpDetails(result.Gcp)); err != nil {
			return err
		}
	}

	return nil
}

func flattenStreamingExportAwsDetails(aws streamingexport.StreamingExportAwsDetails) []interface{} {
	m := map[string]interface{}{
		"aws_account_id":       aws.AwsAccountId,
		"delivery_stream_name": aws.DeliveryStreamName,
		"region":               aws.Region,
		"role":                 aws.Role,
	}
	return []interface{}{m}
}

func flattenStreamingExportAzureDetails(azure streamingexport.StreamingExportAzureDetails) []interface{} {
	m := map[string]interface{}{
		"event_hub_connection_string": azure.EventHubConnectionString,
		"event_hub_name":              azure.EventHubName,
	}
	return []interface{}{m}
}

func flattenStreamingExportGcpDetails(gcp streamingexport.StreamingExportGcpDetails) []interface{} {
	m := map[string]interface{}{
		"gcp_project_id":  gcp.GcpProjectId,
		"pubsub_topic_id": gcp.PubsubTopicId,
	}
	return []interface{}{m}
}
