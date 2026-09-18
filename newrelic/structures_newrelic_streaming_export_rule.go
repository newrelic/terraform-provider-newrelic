package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/streamingexport"
)

func expandStreamingExportRule(d *schema.ResourceData) (streamingexport.StreamingExportAwsInput, streamingexport.StreamingExportAzureInput, streamingexport.StreamingExportGcpInput, streamingexport.StreamingExportRuleInput, error) {
	awsInput := streamingexport.StreamingExportAwsInput{}
	azureInput := streamingexport.StreamingExportAzureInput{}
	gcpInput := streamingexport.StreamingExportGcpInput{}

	ruleInput := streamingexport.StreamingExportRuleInput{
		Name: d.Get("name").(string),
		NRQL: streamingexport.NRQL(d.Get("nrql").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		ruleInput.Description = v.(string)
	}

	if v, ok := d.GetOk("payload_compression"); ok {
		ruleInput.PayloadCompression = streamingexport.StreamingExportPayloadCompression(v.(string))
	}

	if v, ok := d.GetOk("aws"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			awsInput = expandStreamingExportRuleAws(items[0].(map[string]interface{}))
		}
	}

	if v, ok := d.GetOk("azure"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			azureInput = expandStreamingExportRuleAzure(items[0].(map[string]interface{}))
		}
	}

	if v, ok := d.GetOk("gcp"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			gcpInput = expandStreamingExportRuleGcp(items[0].(map[string]interface{}))
		}
	}

	return awsInput, azureInput, gcpInput, ruleInput, nil
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

func flattenStreamingExportRule(result *streamingexport.StreamingExportRule, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	d.SetId(fmt.Sprintf("%d", result.ID))

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
		if err := d.Set("aws", flattenStreamingExportRuleAws(result.Aws)); err != nil {
			return err
		}
	}

	if result.Azure != (streamingexport.StreamingExportAzureDetails{}) {
		if err := d.Set("azure", flattenStreamingExportRuleAzure(result.Azure, d)); err != nil {
			return err
		}
	}

	if result.Gcp != (streamingexport.StreamingExportGcpDetails{}) {
		if err := d.Set("gcp", flattenStreamingExportRuleGcp(result.Gcp)); err != nil {
			return err
		}
	}

	return nil
}

func flattenStreamingExportRuleAws(aws streamingexport.StreamingExportAwsDetails) []interface{} {
	m := map[string]interface{}{
		"aws_account_id":       aws.AwsAccountId,
		"delivery_stream_name": aws.DeliveryStreamName,
		"region":               aws.Region,
		"role":                 aws.Role,
	}
	return []interface{}{m}
}

func flattenStreamingExportRuleAzure(azure streamingexport.StreamingExportAzureDetails, d *schema.ResourceData) []interface{} {
	m := map[string]interface{}{
		"event_hub_name": azure.EventHubName,
	}

	// Preserve sensitive field from state
	if v, ok := d.GetOk("azure.0.event_hub_connection_string"); ok {
		if azure.EventHubConnectionString != "" {
			m["event_hub_connection_string"] = azure.EventHubConnectionString
		} else {
			m["event_hub_connection_string"] = v.(string)
		}
	} else {
		m["event_hub_connection_string"] = azure.EventHubConnectionString
	}

	return []interface{}{m}
}

func flattenStreamingExportRuleGcp(gcp streamingexport.StreamingExportGcpDetails) []interface{} {
	m := map[string]interface{}{
		"gcp_project_id":  gcp.GcpProjectId,
		"pubsub_topic_id": gcp.PubsubTopicId,
	}
	return []interface{}{m}
}
