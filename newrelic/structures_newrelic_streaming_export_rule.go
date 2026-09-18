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

	var awsParameters streamingexport.StreamingExportAwsInput
	if v, ok := d.GetOk("aws"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			awsParameters = expandStreamingExportRuleAws(items[0].(map[string]interface{}))
		}
	}

	var azureParameters streamingexport.StreamingExportAzureInput
	if v, ok := d.GetOk("azure"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			azureParameters = expandStreamingExportRuleAzure(items[0].(map[string]interface{}))
		}
	}

	var gcpParameters streamingexport.StreamingExportGcpInput
	if v, ok := d.GetOk("gcp"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			gcpParameters = expandStreamingExportRuleGcp(items[0].(map[string]interface{}))
		}
	}

	return awsParameters, azureParameters, gcpParameters, ruleParameters, nil
}

func expandStreamingExportRuleAws(cfg map[string]interface{}) streamingexport.StreamingExportAwsInput {
	input := streamingexport.StreamingExportAwsInput{}

	if v, ok := cfg["aws_account_id"]; ok {
		input.AwsAccountId = v.(int)
	}

	if v, ok := cfg["delivery_stream_name"]; ok {
		input.DeliveryStreamName = v.(string)
	}

	if v, ok := cfg["region"]; ok {
		input.Region = v.(string)
	}

	if v, ok := cfg["role"]; ok {
		input.Role = v.(string)
	}

	return input
}

func expandStreamingExportRuleAzure(cfg map[string]interface{}) streamingexport.StreamingExportAzureInput {
	input := streamingexport.StreamingExportAzureInput{}

	if v, ok := cfg["event_hub_connection_string"]; ok {
		input.EventHubConnectionString = v.(string)
	}

	if v, ok := cfg["event_hub_name"]; ok {
		input.EventHubName = v.(string)
	}

	return input
}

func expandStreamingExportRuleGcp(cfg map[string]interface{}) streamingexport.StreamingExportGcpInput {
	input := streamingexport.StreamingExportGcpInput{}

	if v, ok := cfg["gcp_project_id"]; ok {
		input.GcpProjectId = v.(string)
	}

	if v, ok := cfg["pubsub_topic_id"]; ok {
		input.PubsubTopicId = v.(string)
	}

	return input
}

func flattenStreamingExportRule(rule *streamingexport.StreamingExportRule, d *schema.ResourceData) error {
	if rule == nil {
		return nil
	}

	if err := d.Set("name", rule.Name); err != nil {
		return err
	}

	if err := d.Set("nrql", string(rule.NRQL)); err != nil {
		return err
	}

	if err := d.Set("description", rule.Description); err != nil {
		return err
	}

	if err := d.Set("payload_compression", string(rule.PayloadCompression)); err != nil {
		return err
	}

	if err := d.Set("status", string(rule.Status)); err != nil {
		return err
	}

	if err := d.Set("message", rule.Message); err != nil {
		return err
	}

	if err := d.Set("created_at", string(rule.CreatedAt)); err != nil {
		return err
	}

	if err := d.Set("updated_at", string(rule.UpdatedAt)); err != nil {
		return err
	}

	if err := d.Set("account_id", rule.Account.ID); err != nil {
		return err
	}

	if rule.Aws != (streamingexport.StreamingExportAwsDetails{}) {
		if err := d.Set("aws", flattenStreamingExportAwsDetails(rule.Aws)); err != nil {
			return err
		}
	}

	if rule.Azure != (streamingexport.StreamingExportAzureDetails{}) {
		if err := d.Set("azure", flattenStreamingExportAzureDetails(rule.Azure)); err != nil {
			return err
		}
	}

	if rule.Gcp != (streamingexport.StreamingExportGcpDetails{}) {
		if err := d.Set("gcp", flattenStreamingExportGcpDetails(rule.Gcp)); err != nil {
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
