package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewrelicCloudGcpIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewrelicCloudGcpIntegrationsCreate,
		ReadContext:   resourceNewrelicCloudGcpIntegrationsRead,
		UpdateContext: resourceNewrelicCloudGcpIntegrationsUpdate,
		DeleteContext: resourceNewrelicCloudGcpIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: generateGcpIntegrationSchema(),
	}
}

func generateGcpIntegrationSchema() map[string]*schema.Schema {
	baseSchema := cloudGcpIntegrationSchemaBase()
	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Description: "ID of the newrelic account",
			Computed:    true,
			Optional:    true,
		},
		"linked_account_id": {
			Type:        schema.TypeInt,
			Description: "Id of the linked gcp account in New Relic",
			Required:    true,
			ForceNew:    true,
		},
		"alloy_db": {
			Type:        schema.TypeList,
			Description: "GCP alloy DB integration",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"app_engine": {
			Type:        schema.TypeList,
			Description: "GCP app engine service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"big_query": {
			Type:        schema.TypeList,
			Description: "GCP big query service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: addFetchTagsToSchema(baseSchema)},
			MaxItems:    1,
		},
		"big_table": {
			Type:        schema.TypeList,
			Description: "GCP big table service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"composer": {
			Type:        schema.TypeList,
			Description: "GCP composer service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"data_flow": {
			Type:        schema.TypeList,
			Description: "GCP data flow service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"data_proc": {
			Type:        schema.TypeList,
			Description: "GCP data proc service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"data_store": {
			Type:        schema.TypeList,
			Description: "GCP data store service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"fire_base_database": {
			Type:        schema.TypeList,
			Description: "GCP firebase database service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"fire_base_hosting": {
			Type:        schema.TypeList,
			Description: "GCP firebase hosting service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"fire_base_storage": {
			Type:        schema.TypeList,
			Description: "GCP firebase storage service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"fire_store": {
			Type:        schema.TypeList,
			Description: "GCP firestore service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"functions": {
			Type:        schema.TypeList,
			Description: "GCP functions service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"interconnect": {
			Type:        schema.TypeList,
			Description: "GCP interconnect service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"kubernetes": {
			Type:        schema.TypeList,
			Description: "GCP kubernetes service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"load_balancing": {
			Type:        schema.TypeList,
			Description: "GCP load balancing service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"mem_cache": {
			Type:        schema.TypeList,
			Description: "GCP mem cache service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"pub_sub": {
			Type:        schema.TypeList,
			Description: "GCP pub sub service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: addFetchTagsToSchema(baseSchema)},
			MaxItems:    1,
		},
		"redis": {
			Type:        schema.TypeList,
			Description: "GCP redis service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"router": {
			Type:        schema.TypeList,
			Description: "GCP router service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"run": {
			Type:        schema.TypeList,
			Description: "GCP run service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"spanner": {
			Type:        schema.TypeList,
			Description: "GCP spanner service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: addFetchTagsToSchema(baseSchema)},
			MaxItems:    1,
		},
		"sql": {
			Type:        schema.TypeList,
			Description: "GCP SQL service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"storage": {
			Type:        schema.TypeList,
			Description: "GCP storage service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: addFetchTagsToSchema(baseSchema)},
			MaxItems:    1,
		},
		"virtual_machines": {
			Type:        schema.TypeList,
			Description: "GCP virtual machines service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
		"vpc_access": {
			Type:        schema.TypeList,
			Description: "GCP VPC access service",
			Optional:    true,
			Elem:        &schema.Resource{Schema: baseSchema},
			MaxItems:    1,
		},
	}
}

func addFetchTagsToSchema(baseSchema map[string]*schema.Schema) map[string]*schema.Schema {
	schemaWithTags := make(map[string]*schema.Schema)
	for k, v := range baseSchema {
		schemaWithTags[k] = v
	}
	schemaWithTags["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "to fetch tags of the resource",
		Optional:    true,
	}
	return schemaWithTags
}

func resourceNewrelicCloudGcpIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	cloudGcpIntegrationinputs, _ := expandCloudGcpIntegrationsinputs(d)
	gcpIntegrationspayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudGcpIntegrationinputs)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if len(gcpIntegrationspayload.Errors) > 0 {
		for _, err := range gcpIntegrationspayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	if len(gcpIntegrationspayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}
	return nil
}

// expand function to extract inputs for cloud integrations from the schema
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandCloudGcpIntegrationsinputs(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	gcpCloudIntegrations := cloud.CloudGcpIntegrationsInput{}
	gcpDisableIntegrations := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if lid, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = lid.(int)
	}
	if v, ok := d.GetOk("alloy_db"); ok {
		gcpCloudIntegrations.GcpAlloydb = expandCloudGcpAlloyDBIntegrationsInputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("alloy_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("app_engine"); ok {
		gcpCloudIntegrations.GcpAppengine = expandCloudGcpAppEngineIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("app_engine"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("big_query"); ok {
		gcpCloudIntegrations.GcpBigquery = expandCloudGcpBigQueryIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("big_query"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("big_table"); ok {
		gcpCloudIntegrations.GcpBigtable = expandCloudGcpBigTableIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("big_table"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("composer"); ok {
		gcpCloudIntegrations.GcpComposer = expandCloudGcpComposerIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("composer"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpComposer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("data_flow"); ok {
		gcpCloudIntegrations.GcpDataflow = expandCloudGcpDataFlowIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("data_flow"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpDataflow = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("data_proc"); ok {
		gcpCloudIntegrations.GcpDataproc = expandCloudGcpDataProcIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("data_proc"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpDataproc = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("data_store"); ok {
		gcpCloudIntegrations.GcpDatastore = expandCloudGcpDataStoreIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("data_store"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpDatastore = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("fire_base_database"); ok {
		gcpCloudIntegrations.GcpFirebasedatabase = expandCloudGcpFireBaseDatabaseIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("fire_base_database"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpFirebasedatabase = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("fire_base_hosting"); ok {
		gcpCloudIntegrations.GcpFirebasehosting = expandCloudGcpFireBaseHostingIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("fire_base_hosting"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpFirebasehosting = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("fire_base_storage"); ok {
		gcpCloudIntegrations.GcpFirebasestorage = expandCloudGcpFireBaseStorageIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("fire_base_storage"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpFirebasestorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("fire_store"); ok {
		gcpCloudIntegrations.GcpFirestore = expandCloudGcpFireStoreIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("fire_store"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpFirestore = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("functions"); ok {
		gcpCloudIntegrations.GcpFunctions = expandCloudGcpFunctionsIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("functions"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("interconnect"); ok {
		gcpCloudIntegrations.GcpInterconnect = expandCloudGcpInterconnectIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("interconnect"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpInterconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("kubernetes"); ok {
		gcpCloudIntegrations.GcpKubernetes = expandCloudGcpKubernetesIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("kubernetes"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpKubernetes = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("load_balancing"); ok {
		gcpCloudIntegrations.GcpLoadbalancing = expandCloudGcpLoadBalancingIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("load_balancing"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpLoadbalancing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("mem_cache"); ok {
		gcpCloudIntegrations.GcpMemcache = expandCloudGcpMemCacheIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("mem_cache"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpMemcache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("pub_sub"); ok {
		gcpCloudIntegrations.GcpPubsub = expandCloudGcpPubSubIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("pub_sub"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("redis"); ok {
		gcpCloudIntegrations.GcpRedis = expandCloudGcpRedisIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("redis"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpRedis = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("router"); ok {
		gcpCloudIntegrations.GcpRouter = expandCloudGcpRouterIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("router"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpRouter = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("run"); ok {
		gcpCloudIntegrations.GcpRun = expandCloudGcpRunIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("run"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpRun = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("spanner"); ok {
		gcpCloudIntegrations.GcpSpanner = expandCloudGcpSpannerIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("spanner"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("sql"); ok {
		gcpCloudIntegrations.GcpSql = expandCloudGcpSQLIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("sql"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("storage"); ok {
		gcpCloudIntegrations.GcpStorage = expandCloudGcpStorageIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("storage"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("virtual_machines"); ok {
		gcpCloudIntegrations.GcpVms = expandCloudGcpVirtualMachinesIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("virtual_machines"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("vpc_access"); ok {
		gcpCloudIntegrations.GcpVpcaccess = expandCloudGcpVpcAccessIntegrationsinputs(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("vpc_access"); len(n.([]interface{})) < len(o.([]interface{})) {
		gcpDisableIntegrations.GcpVpcaccess = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	configureInput := cloud.CloudIntegrationsInput{
		Gcp: gcpCloudIntegrations,
	}
	disableInput := cloud.CloudDisableIntegrationsInput{
		Gcp: gcpDisableIntegrations,
	}
	return configureInput, disableInput
}

// expand function to extract inputs from gcp app engine schema
func expandCloudGcpAppEngineIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpAppengineIntegrationInput {
	expanded := make([]cloud.CloudGcpAppengineIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpAppengineIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if a, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = a.(int)
		}
		expanded[i] = input
	}
	return expanded
}

func expandCloudGcpAlloyDBIntegrationsInputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpAlloydbIntegrationInput {
	expanded := make([]cloud.CloudGcpAlloydbIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpAlloydbIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if a, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = a.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp bigquery schema
func expandCloudGcpBigQueryIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpBigqueryIntegrationInput {
	expanded := make([]cloud.CloudGcpBigqueryIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpBigqueryIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		if f, ok := in["fetch_tags"]; ok {
			input.FetchTags = f.(bool)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp bigtable schema
func expandCloudGcpBigTableIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpBigtableIntegrationInput {
	expanded := make([]cloud.CloudGcpBigtableIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpBigtableIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp composer schema
func expandCloudGcpComposerIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpComposerIntegrationInput {
	expanded := make([]cloud.CloudGcpComposerIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpComposerIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp dataflow schema
func expandCloudGcpDataFlowIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpDataflowIntegrationInput {
	expanded := make([]cloud.CloudGcpDataflowIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpDataflowIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp dataproc schema
func expandCloudGcpDataProcIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpDataprocIntegrationInput {
	expanded := make([]cloud.CloudGcpDataprocIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpDataprocIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp datastore schema
func expandCloudGcpDataStoreIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpDatastoreIntegrationInput {
	expanded := make([]cloud.CloudGcpDatastoreIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpDatastoreIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp firebasedatabase schema
func expandCloudGcpFireBaseDatabaseIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpFirebasedatabaseIntegrationInput {
	expanded := make([]cloud.CloudGcpFirebasedatabaseIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpFirebasedatabaseIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp firebasehosting schema
func expandCloudGcpFireBaseHostingIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpFirebasehostingIntegrationInput {
	expanded := make([]cloud.CloudGcpFirebasehostingIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpFirebasehostingIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp firebasestorage schema
func expandCloudGcpFireBaseStorageIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpFirebasestorageIntegrationInput {
	expanded := make([]cloud.CloudGcpFirebasestorageIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpFirebasestorageIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp firestore schema
func expandCloudGcpFireStoreIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpFirestoreIntegrationInput {
	expanded := make([]cloud.CloudGcpFirestoreIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpFirestoreIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp functions schema
func expandCloudGcpFunctionsIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpFunctionsIntegrationInput {
	expanded := make([]cloud.CloudGcpFunctionsIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpFunctionsIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp interconnect schema
func expandCloudGcpInterconnectIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpInterconnectIntegrationInput {
	expanded := make([]cloud.CloudGcpInterconnectIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpInterconnectIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp k8's schema
func expandCloudGcpKubernetesIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpKubernetesIntegrationInput {
	expanded := make([]cloud.CloudGcpKubernetesIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpKubernetesIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp load balancing schema
func expandCloudGcpLoadBalancingIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpLoadbalancingIntegrationInput {
	expanded := make([]cloud.CloudGcpLoadbalancingIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpLoadbalancingIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp mem-cache schema
func expandCloudGcpMemCacheIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpMemcacheIntegrationInput {
	expanded := make([]cloud.CloudGcpMemcacheIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpMemcacheIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp pubsub schema
func expandCloudGcpPubSubIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpPubsubIntegrationInput {
	expanded := make([]cloud.CloudGcpPubsubIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpPubsubIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		if f, ok := in["fetch_tags"]; ok {
			input.FetchTags = f.(bool)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp redis schema
func expandCloudGcpRedisIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpRedisIntegrationInput {
	expanded := make([]cloud.CloudGcpRedisIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpRedisIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp router schema
func expandCloudGcpRouterIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpRouterIntegrationInput {
	expanded := make([]cloud.CloudGcpRouterIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpRouterIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp run schema
func expandCloudGcpRunIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpRunIntegrationInput {
	expanded := make([]cloud.CloudGcpRunIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpRunIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp spanner schema
func expandCloudGcpSpannerIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpSpannerIntegrationInput {
	expanded := make([]cloud.CloudGcpSpannerIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpSpannerIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		if f, ok := in["fetch_tags"]; ok {
			input.FetchTags = f.(bool)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp SQL schema
func expandCloudGcpSQLIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpSqlIntegrationInput {
	expanded := make([]cloud.CloudGcpSqlIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpSqlIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp storage schema
func expandCloudGcpStorageIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpStorageIntegrationInput {
	expanded := make([]cloud.CloudGcpStorageIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpStorageIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		if f, ok := in["fetch_tags"]; ok {
			input.FetchTags = f.(bool)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp vm's schema
func expandCloudGcpVirtualMachinesIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpVmsIntegrationInput {
	expanded := make([]cloud.CloudGcpVmsIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpVmsIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

// expand function to extract inputs from gcp vpc access schema
func expandCloudGcpVpcAccessIntegrationsinputs(b []interface{}, linkedAccountID int) []cloud.CloudGcpVpcaccessIntegrationInput {
	expanded := make([]cloud.CloudGcpVpcaccessIntegrationInput, len(b))
	for i, expand := range b {
		var input cloud.CloudGcpVpcaccessIntegrationInput
		if expand == nil {
			input.LinkedAccountId = linkedAccountID
			expanded[i] = input
			return expanded
		}
		in := expand.(map[string]interface{})
		input.LinkedAccountId = linkedAccountID
		if m, ok := in["metrics_polling_interval"]; ok {
			input.MetricsPollingInterval = m.(int)
		}
		expanded[i] = input
	}
	return expanded
}

func resourceNewrelicCloudGcpIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	flattenCloudGcpLinkedAccount(d, linkedAccount)
	return nil
}

func cloudGcpIntegrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_engine": cloud.CloudGcpAppengineIntegration,
	}
}

// flatten function to set(store) outputs from the terraform apply
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func flattenCloudGcpLinkedAccount(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)
	for _, i := range linkedAccount.Integrations {
		switch t := i.(type) {
		case *cloud.CloudGcpAlloydbIntegration:
			_ = d.Set("alloy_db", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpAppengineIntegration:
			_ = d.Set("app_engine", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpBigqueryIntegration:
			_ = d.Set("big_query", flattenCloudGcpBigQueryIntegration(t))
		case *cloud.CloudGcpBigtableIntegration:
			_ = d.Set("big_table", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpComposerIntegration:
			_ = d.Set("composer", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpDataflowIntegration:
			_ = d.Set("data_flow", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpDataprocIntegration:
			_ = d.Set("data_proc", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpDatastoreIntegration:
			_ = d.Set("data_store", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpFirebasedatabaseIntegration:
			_ = d.Set("fire_base_database", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpFirebasehostingIntegration:
			_ = d.Set("fire_base_hosting", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpFirebasestorageIntegration:
			_ = d.Set("fire_base_storage", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpFirestoreIntegration:
			_ = d.Set("fire_store", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpFunctionsIntegration:
			_ = d.Set("functions", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpInterconnectIntegration:
			_ = d.Set("interconnect", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpKubernetesIntegration:
			_ = d.Set("kubernetes", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpLoadbalancingIntegration:
			_ = d.Set("load_balancing", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpMemcacheIntegration:
			_ = d.Set("mem_cache", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpPubsubIntegration:
			_ = d.Set("pub_sub", flattenCloudGcpBigQueryIntegration(t))
		case *cloud.CloudGcpRedisIntegration:
			_ = d.Set("redis", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpRouterIntegration:
			_ = d.Set("router", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpRunIntegration:
			_ = d.Set("run", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpSpannerIntegration:
			_ = d.Set("spanner", flattenCloudGcpBigQueryIntegration(t))
		case *cloud.CloudGcpSqlIntegration:
			_ = d.Set("sql", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpStorageIntegration:
			_ = d.Set("storage", flattenCloudGcpBigQueryIntegration(t))
		case *cloud.CloudGcpVmsIntegration:
			_ = d.Set("virtual_machines", flattenCloudGcpCommonIntegration(t))
		case *cloud.CloudGcpVpcaccessIntegration:
			_ = d.Set("vpc_access", flattenCloudGcpCommonIntegration(t))
		}
	}
}

func flattenCloudGcpCommonIntegration(in interface{}) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})

	switch t := in.(type) {
	case *cloud.CloudGcpBigqueryIntegration, *cloud.CloudGcpPubsubIntegration, *cloud.CloudGcpSpannerIntegration, *cloud.CloudGcpStorageIntegration:
		out["fetch_tags"] = t.FetchTags
	}
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

func resourceNewrelicCloudGcpIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	configureInput, disableInput := expandCloudGcpIntegrationsinputs(d)
	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudDisableIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudDisableIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	cloudGcpIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(cloudGcpIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudGcpIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

func resourceNewrelicCloudGcpIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	deleteInput := expandCloudGcpDisableinputs(d)
	gcpDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(gcpDisablePayload.Errors) > 0 {
		for _, err := range gcpDisablePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	d.SetId("")

	return nil
}

// expand function to extract the inputs values from the schema for disabling the integration for particular services
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandCloudGcpDisableinputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudGcpDisableInput := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("alloy_db"); ok {
		cloudGcpDisableInput.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("app_engine"); ok {
		cloudGcpDisableInput.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("big_query"); ok {
		cloudGcpDisableInput.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("big_table"); ok {
		cloudGcpDisableInput.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("composer"); ok {
		cloudGcpDisableInput.GcpComposer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("data_flow"); ok {
		cloudGcpDisableInput.GcpDataflow = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("data_proc"); ok {
		cloudGcpDisableInput.GcpDataproc = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("data_store"); ok {
		cloudGcpDisableInput.GcpDatastore = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("fire_base_database"); ok {
		cloudGcpDisableInput.GcpFirebasedatabase = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("fire_base_hosting"); ok {
		cloudGcpDisableInput.GcpFirebasehosting = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("fire_base_storage"); ok {
		cloudGcpDisableInput.GcpFirebasestorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("fire_store"); ok {
		cloudGcpDisableInput.GcpFirestore = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("functions"); ok {
		cloudGcpDisableInput.GcpFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("interconnect"); ok {
		cloudGcpDisableInput.GcpInterconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("kubernetes"); ok {
		cloudGcpDisableInput.GcpKubernetes = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("load_balancing"); ok {
		cloudGcpDisableInput.GcpLoadbalancing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("mem_cache"); ok {
		cloudGcpDisableInput.GcpMemcache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("pub_sub"); ok {
		cloudGcpDisableInput.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("redis"); ok {
		cloudGcpDisableInput.GcpRedis = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("router"); ok {
		cloudGcpDisableInput.GcpRouter = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("run"); ok {
		cloudGcpDisableInput.GcpRun = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("spanner"); ok {
		cloudGcpDisableInput.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("sql"); ok {
		cloudGcpDisableInput.GcpSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("storage"); ok {
		cloudGcpDisableInput.GcpStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("virtual_machines"); ok {
		cloudGcpDisableInput.GcpVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("vpc_access"); ok {
		cloudGcpDisableInput.GcpVpcaccess = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		Gcp: cloudGcpDisableInput,
	}
	return deleteInput
}
