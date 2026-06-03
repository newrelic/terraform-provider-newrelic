package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/federatedlogs"
)

// Recreate hints used by the setup / partition CustomizeDiff validators when
// reporting that a field can't be updated in place. Centralised here so the
// per-field error sites stay one-liners.
const (
	federatedLogsSetupRecreateHint     = "recreate the setup to change it (this destroys the default partition and every sub-partition under it)"
	federatedLogsPartitionRecreateHint = "recreate the partition to change it"
)

// immutableFieldError builds the standard error returned when a customer
// tries to edit a create-only field on an existing federated logs resource.
// `field` is the human-readable attribute path (no TypeList ".0." indexing);
// `hint` describes how to actually apply the change.
func immutableFieldError(field, hint string) error {
	return fmt.Errorf("%s cannot be updated after creation; %s", field, hint)
}

// =====================================================================
// Shared schema helpers
// =====================================================================

// statusDetailSchema returns the {status, message, last_updated_at} block
// shape used both for lifecycle_status (where status is a state-machine value
// like COMPLETE / DELETING / ERROR) and for each entry under health_check
// (where status is a probe outcome like HEALTHY / UNHEALTHY / UNKNOWN). The
// underlying types differ but the schema fields are identical, so they share
// this helper. All fields are Computed because the API populates them and the
// terraform user does not author them directly.
func statusDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status":          {Type: schema.TypeString, Computed: true},
			"message":         {Type: schema.TypeString, Computed: true},
			"last_updated_at": {Type: schema.TypeString, Computed: true},
		},
	}
}

// =====================================================================
// Setup expand / flatten helpers
// =====================================================================

func expandFederatedLogsSetupStorage(in []interface{}) federatedlogs.FederatedLogsSetupStorageInput {
	out := federatedlogs.FederatedLogsSetupStorageInput{}
	if len(in) == 0 || in[0] == nil {
		return out
	}
	m := in[0].(map[string]interface{})
	out.DataLocationBucket = m["data_location_bucket"].(string)
	out.Database = m["database"].(string)
	out.DataIngestConnectionId = m["data_ingest_connection_id"].(string)
	out.QueryConnectionId = m["query_connection_id"].(string)
	out.CloudProviderConfiguration = expandFederatedLogsCloudProviderConfig(m["cloud_provider_configuration"].([]interface{}))
	return out
}

func expandFederatedLogsCloudProviderConfig(in []interface{}) federatedlogs.FederatedLogsCloudProviderConfigurationInput {
	out := federatedlogs.FederatedLogsCloudProviderConfigurationInput{}
	if len(in) == 0 || in[0] == nil {
		return out
	}
	m := in[0].(map[string]interface{})
	out.Provider = federatedlogs.FederatedLogsCloudProvider(m["provider"].(string))
	out.Region = m["region"].(string)
	return out
}

func expandFederatedLogsDefaultPartition(in []interface{}) federatedlogs.FederatedLogsDefaultPartitionInput {
	out := federatedlogs.FederatedLogsDefaultPartitionInput{}
	if len(in) == 0 || in[0] == nil {
		return out
	}
	m := in[0].(map[string]interface{})
	out.Storage = expandFederatedLogsPartitionStorage(m["storage"].([]interface{}))
	if v, ok := m["data_retention_policy"]; ok {
		out.DataRetentionPolicy = expandFederatedLogsRetentionPolicy(v.([]interface{}))
	}
	return out
}

func expandFederatedLogsDefaultPartitionUpdate(in []interface{}) *federatedlogs.FederatedLogsUpdateDefaultPartitionInput {
	policy := expandFederatedLogsRetentionPolicy(in)
	if policy == nil {
		return nil
	}
	return &federatedlogs.FederatedLogsUpdateDefaultPartitionInput{
		DataRetentionPolicy: *policy,
	}
}

func expandFederatedLogsForwarder(in []interface{}) *federatedlogs.FederatedLogsForwarderInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	out := federatedlogs.FederatedLogsForwarderInput{
		Type: federatedlogs.FederatedLogsForwarderType(m["type"].(string)),
	}
	if v, ok := m["pipeline_control"]; ok {
		out.PipelineControl = expandFederatedLogsPipelineControl(v.([]interface{}))
	}
	return &out
}

func expandFederatedLogsPipelineControl(in []interface{}) *federatedlogs.FederatedLogsPipelineControlConfigurationInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	out := federatedlogs.FederatedLogsPipelineControlConfigurationInput{
		FleetId: m["fleet_id"].(string),
	}
	if v, ok := m["routing_rule"]; ok {
		out.RoutingRule = expandFederatedLogsRule(v.([]interface{}))
	}
	return &out
}

func expandFederatedLogsRule(in []interface{}) *federatedlogs.FederatedLogsRuleInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	return &federatedlogs.FederatedLogsRuleInput{
		Expression: m["expression"].(string),
	}
}

// =====================================================================
// Partition expand helpers
// =====================================================================

func expandFederatedLogsPartitionStorage(in []interface{}) federatedlogs.FederatedLogsPartitionStorageInput {
	out := federatedlogs.FederatedLogsPartitionStorageInput{}
	if len(in) == 0 || in[0] == nil {
		return out
	}
	m := in[0].(map[string]interface{})
	out.Table = m["table"].(string)
	out.DataLocationUri = m["data_location_uri"].(string)
	return out
}

func expandFederatedLogsRetentionPolicy(in []interface{}) *federatedlogs.FederatedLogsRetentionPolicyInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	return &federatedlogs.FederatedLogsRetentionPolicyInput{
		Duration: m["duration"].(int),
		Unit:     federatedlogs.FederatedLogsRetentionUnit(m["unit"].(string)),
	}
}

func expandFederatedLogsPartitionForwarderConfig(in []interface{}) *federatedlogs.FederatedLogsPartitionForwarderConfigurationInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	out := federatedlogs.FederatedLogsPartitionForwarderConfigurationInput{
		Type: federatedlogs.FederatedLogsForwarderType(m["type"].(string)),
	}
	if v, ok := m["pipeline_control"]; ok {
		out.PipelineControl = expandFederatedLogsPartitionPipelineControl(v.([]interface{}))
	}
	return &out
}

func expandFederatedLogsPartitionPipelineControl(in []interface{}) *federatedlogs.FederatedLogsPartitionPipelineControlConfigurationInput {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m := in[0].(map[string]interface{})
	out := federatedlogs.FederatedLogsPartitionPipelineControlConfigurationInput{}
	if v, ok := m["partition_rule"]; ok {
		out.PartitionRule = expandFederatedLogsRule(v.([]interface{}))
	}
	return &out
}

// =====================================================================
// Setup flatten (read response → state)
// =====================================================================

func flattenFederatedLogsSetupIntoState(d *schema.ResourceData, s *federatedlogs.FederatedLogsSetup) error {
	if err := d.Set("name", s.Name); err != nil {
		return err
	}
	if err := d.Set("description", s.Description); err != nil {
		return err
	}
	if err := d.Set("active", s.Active); err != nil {
		return err
	}
	if err := d.Set("default_partition_id", s.DefaultPartitionId); err != nil {
		return err
	}
	if err := d.Set("storage", flattenFederatedLogsSetupStorage(s.Storage)); err != nil {
		return err
	}
	if err := d.Set("forwarder", flattenFederatedLogsForwarder(s.Forwarder)); err != nil {
		return err
	}
	if err := d.Set("lifecycle_status", []map[string]interface{}{{
		"status":          string(s.LifecycleStatus.Status),
		"message":         s.LifecycleStatus.Message,
		"last_updated_at": string(s.LifecycleStatus.LastUpdatedAt),
	}}); err != nil {
		return err
	}
	if err := d.Set("health_check", flattenFederatedLogsSetupHealthCheck(s.HealthCheck)); err != nil {
		return err
	}
	if err := d.Set("created_at", string(s.CreatedAt)); err != nil {
		return err
	}
	if err := d.Set("updated_at", string(s.UpdatedAt)); err != nil {
		return err
	}
	return nil
}

func flattenFederatedLogsSetupStorage(s federatedlogs.FederatedLogsSetupStorage) []map[string]interface{} {
	return []map[string]interface{}{{
		"data_location_bucket":      s.DataLocationBucket,
		"database":                  s.Database,
		"data_ingest_connection_id": s.DataIngestConnectionId,
		"query_connection_id":       s.QueryConnectionId,
		"cloud_provider_configuration": []map[string]interface{}{{
			"provider": string(s.CloudProviderConfiguration.Provider),
			"region":   s.CloudProviderConfiguration.Region,
		}},
	}}
}

func flattenFederatedLogsForwarder(f federatedlogs.FederatedLogsForwarder) []map[string]interface{} {
	if f.Type == "" {
		return nil
	}
	out := map[string]interface{}{
		"type": string(f.Type),
	}
	if f.PipelineControl.FleetId != "" {
		pc := map[string]interface{}{
			"fleet_id": f.PipelineControl.FleetId,
		}
		if f.PipelineControl.RoutingRule.Expression != "" {
			pc["routing_rule"] = []map[string]interface{}{{"expression": f.PipelineControl.RoutingRule.Expression}}
		}
		out["pipeline_control"] = []map[string]interface{}{pc}
	}
	return []map[string]interface{}{out}
}

func flattenFederatedLogsSetupHealthCheck(h federatedlogs.FederatedLogsSetupHealthCheckStatus) []map[string]interface{} {
	out := map[string]interface{}{
		"last_updated_at":   string(h.LastUpdatedAt),
		"query_connection":  flattenFederatedLogsHealthCheckDetail(h.QueryConnection),
		"end2end_data_flow": flattenFederatedLogsHealthCheckDetail(h.End2endDataFlow),
	}
	return []map[string]interface{}{out}
}

func flattenFederatedLogsHealthCheckDetail(h federatedlogs.FederatedLogsHealthCheckDetail) []map[string]interface{} {
	if h.Status == "" {
		return nil
	}
	return []map[string]interface{}{{
		"status":          string(h.Status),
		"message":         h.Message,
		"last_updated_at": string(h.LastUpdatedAt),
	}}
}

// =====================================================================
// Partition flatten (read response → state)
// =====================================================================

func flattenFederatedLogsPartitionIntoState(d *schema.ResourceData, p *federatedlogs.FederatedLogsPartition) error {
	if err := d.Set("setup_id", p.Setup.ID); err != nil {
		return err
	}
	if err := d.Set("name", p.Name); err != nil {
		return err
	}
	if err := d.Set("description", p.Description); err != nil {
		return err
	}
	if err := d.Set("active", p.Active); err != nil {
		return err
	}
	if err := d.Set("is_default", p.IsDefault); err != nil {
		return err
	}
	if err := d.Set("storage", flattenFederatedLogsPartitionStorage(p.Storage)); err != nil {
		return err
	}
	if err := d.Set("data_retention_policy", flattenFederatedLogsRetentionPolicy(p.DataRetentionPolicy)); err != nil {
		return err
	}
	if err := d.Set("forwarder_configuration", flattenFederatedLogsPartitionForwarderConfig(p.ForwarderConfiguration)); err != nil {
		return err
	}
	if err := d.Set("lifecycle_status", []map[string]interface{}{{
		"status":          string(p.LifecycleStatus.Status),
		"message":         p.LifecycleStatus.Message,
		"last_updated_at": string(p.LifecycleStatus.LastUpdatedAt),
	}}); err != nil {
		return err
	}
	if err := d.Set("health_check", flattenFederatedLogsPartitionHealthCheck(p.HealthCheck)); err != nil {
		return err
	}
	if err := d.Set("created_at", string(p.CreatedAt)); err != nil {
		return err
	}
	if err := d.Set("updated_at", string(p.UpdatedAt)); err != nil {
		return err
	}
	return nil
}

func flattenFederatedLogsPartitionStorage(s federatedlogs.FederatedLogsPartitionStorage) []map[string]interface{} {
	return []map[string]interface{}{{
		"table":             s.Table,
		"data_location_uri": s.DataLocationUri,
	}}
}

func flattenFederatedLogsRetentionPolicy(p federatedlogs.FederatedLogsRetentionPolicy) []map[string]interface{} {
	if p.Unit == "" {
		return nil
	}
	return []map[string]interface{}{{
		"duration": p.Duration,
		"unit":     string(p.Unit),
	}}
}

func flattenFederatedLogsPartitionForwarderConfig(c federatedlogs.FederatedLogsPartitionForwarderConfiguration) []map[string]interface{} {
	if c.Type == "" {
		return nil
	}
	out := map[string]interface{}{
		"type": string(c.Type),
	}
	if c.PipelineControl.PartitionRule.Expression != "" {
		out["pipeline_control"] = []map[string]interface{}{{
			"partition_rule": []map[string]interface{}{{
				"expression": c.PipelineControl.PartitionRule.Expression,
			}},
		}}
	}
	return []map[string]interface{}{out}
}

func flattenFederatedLogsPartitionHealthCheck(h federatedlogs.FederatedLogsPartitionHealthCheckStatus) []map[string]interface{} {
	out := map[string]interface{}{
		"last_updated_at":   string(h.LastUpdatedAt),
		"query_connection":  flattenFederatedLogsHealthCheckDetail(h.QueryConnection),
		"end2end_data_flow": flattenFederatedLogsHealthCheckDetail(h.End2endDataFlow),
	}
	return []map[string]interface{}{out}
}
