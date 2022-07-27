package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
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
		Schema: map[string]*schema.Schema{
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
			},
			"app_engine": {
				Type:        schema.TypeList,
				Description: "GCP app engine service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsAppEngineSchemaElem(),
				MaxItems:    1,
			},
			"big_query": {
				Type:        schema.TypeList,
				Description: "GCP biq query service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsBigQuerySchemaElem(),
				MaxItems:    1,
			},
			"big_table": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsBigTableSchemaElem(),
				MaxItems:    1,
			},
			"composer": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsComposerSchemaElem(),
				MaxItems:    1,
			},
			"data_flow": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsDataFlowSchemaElem(),
				MaxItems:    1,
			},
			"data_proc": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsDataProcSchemaElem(),
				MaxItems:    1,
			},
			"data_store": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsDataStoreSchemaElem(),
				MaxItems:    1,
			},
			"fire_base_database": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsFireBaseDatabaseSchemaElem(),
				MaxItems:    1,
			},
			"fire_base_hosting": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsFireBaseHostingSchemaElem(),
				MaxItems:    1,
			},
			"fire_base_storage": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsFireBaseStorageSchemaElem(),
				MaxItems:    1,
			},
			"fire_store": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsFireStoreSchemaElem(),
				MaxItems:    1,
			},
			"functions": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsFunctionsSchemaElem(),
				MaxItems:    1,
			},
			"interconnect": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsInterconnectSchemaElem(),
				MaxItems:    1,
			},
			"kubernetes": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsKubernetesSchemaElem(),
				MaxItems:    1,
			},
			"load_balancing": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsLoadBalancingSchemaElem(),
				MaxItems:    1,
			},
			"mem_cache": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsMemCacheSchemaElem(),
				MaxItems:    1,
			},
			"pub_sub": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsPubSubSchemaElem(),
				MaxItems:    1,
			},
			"redis": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsRedisSchemaElem(),
				MaxItems:    1,
			},
			"router": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsRouterSchemaElem(),
				MaxItems:    1,
			},
			"run": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsRunSchemaElem(),
				MaxItems:    1,
			},
			"spanner": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsSpannerSchemaElem(),
				MaxItems:    1,
			},
			"sql": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsSQLSchemaElem(),
				MaxItems:    1,
			},
			"storage": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsStorageSchemaElem(),
				MaxItems:    1,
			},
			"virtual_machines": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsVirtualMachinesSchemaElem(),
				MaxItems:    1,
			},
			"vpc_access": {
				Type:        schema.TypeList,
				Description: "GCP big table service",
				Optional:    true,
				Elem:        cloudGcpIntegrationsVpcAccessSchemaElem(),
				MaxItems:    1,
			},
		},
	}
}

//function to add common schema for gcp all resources
func cloudGcpIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Description: "the data polling interval in seconds",
			Optional:    true,
		},
	}
}

//function to add schema for gcp AppEngine
func cloudGcpIntegrationsAppEngineSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp bigquery
func cloudGcpIntegrationsBigQuerySchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "to fetch tags of the resource",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp bigtable
func cloudGcpIntegrationsBigTableSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp composer
func cloudGcpIntegrationsComposerSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp dataflow
func cloudGcpIntegrationsDataFlowSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp dataproc
func cloudGcpIntegrationsDataProcSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp datastore
func cloudGcpIntegrationsDataStoreSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp firebasedatabase
func cloudGcpIntegrationsFireBaseDatabaseSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp firebasehosting
func cloudGcpIntegrationsFireBaseHostingSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp firebasestorage
func cloudGcpIntegrationsFireBaseStorageSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp firestore
func cloudGcpIntegrationsFireStoreSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp functions
func cloudGcpIntegrationsFunctionsSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp interconnect
func cloudGcpIntegrationsInterconnectSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp k8's
func cloudGcpIntegrationsKubernetesSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp load balancing
func cloudGcpIntegrationsLoadBalancingSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp mem-cache
func cloudGcpIntegrationsMemCacheSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp pubsub
func cloudGcpIntegrationsPubSubSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "to fetch tags of the resource",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp redis
func cloudGcpIntegrationsRedisSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp router
func cloudGcpIntegrationsRouterSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp run schema
func cloudGcpIntegrationsRunSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp spanner
func cloudGcpIntegrationsSpannerSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "to fetch tags of the resource",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp SQL
func cloudGcpIntegrationsSQLSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp storage
func cloudGcpIntegrationsStorageSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "to fetch tags of the resource",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp vm's
func cloudGcpIntegrationsVirtualMachinesSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for gcp vpc access
func cloudGcpIntegrationsVpcAccessSchemaElem() *schema.Resource {
	s := cloudGcpIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

func resourceNewrelicCloudGcpIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	cloudGcpIntegrationinputs, _ := expandCloudGcpIntegrationsinputs(d)
	gcpIntegrationspayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudGcpIntegrationinputs)
	if err != nil {
		diag.FromErr(err)
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

//expand function to extract inputs for cloud integrations from the schema
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandCloudGcpIntegrationsinputs(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	gcpCloudIntegrations := cloud.CloudGcpIntegrationsInput{}
	gcpDisableIntegrations := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if lid, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = lid.(int)
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

//expand function to extract inputs from gcp app engine schema
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

//expand function to extract inputs from gcp bigquery schema
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

//expand function to extract inputs from gcp bigtable schema
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

//expand function to extract inputs from gcp composer schema
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

//expand function to extract inputs from gcp dataflow schema
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

//expand function to extract inputs from gcp dataproc schema
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

//expand function to extract inputs from gcp datastore schema
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

//expand function to extract inputs from gcp firebasedatabase schema
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

//expand function to extract inputs from gcp firebasehosting schema
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

//expand function to extract inputs from gcp firebasestorage schema
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

//expand function to extract inputs from gcp firestore schema
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

//expand function to extract inputs from gcp functions schema
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

//expand function to extract inputs from gcp interconnect schema
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

//expand function to extract inputs from gcp k8's schema
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

//expand function to extract inputs from gcp load balancing schema
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

//expand function to extract inputs from gcp mem-cache schema
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

//expand function to extract inputs from gcp pubsub schema
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

//expand function to extract inputs from gcp redis schema
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

//expand function to extract inputs from gcp router schema
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

//expand function to extract inputs from gcp run schema
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

//expand function to extract inputs from gcp spanner schema
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

//expand function to extract inputs from gcp SQL schema
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

//expand function to extract inputs from gcp storage schema
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

//expand function to extract inputs from gcp vm's schema
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

//expand function to extract inputs from gcp vpc access schema
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
		diag.FromErr(err)
	}
	flattenCloudGcpLinkedAccount(d, linkedAccount)
	return nil
}

//flatten function to set(store) outputs from the terraform apply
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func flattenCloudGcpLinkedAccount(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)
	for _, i := range linkedAccount.Integrations {
		switch t := i.(type) {
		case *cloud.CloudGcpAppengineIntegration:
			_ = d.Set("app_engine", flattenCloudGcpAppEngineIntegration(t))
		case *cloud.CloudGcpBigqueryIntegration:
			_ = d.Set("big_query", flattenCloudGcpBigQueryIntegration(t))
		case *cloud.CloudGcpBigtableIntegration:
			_ = d.Set("big_table", flattenCloudGcpBigTableIntegration(t))
		case *cloud.CloudGcpComposerIntegration:
			_ = d.Set("composer", flattenCloudGcpComposerIntegration(t))
		case *cloud.CloudGcpDataflowIntegration:
			_ = d.Set("data_flow", flattenCloudGcpDataFlowIntegration(t))
		case *cloud.CloudGcpDataprocIntegration:
			_ = d.Set("data_proc", flattenCloudGcpDataProcIntegration(t))
		case *cloud.CloudGcpDatastoreIntegration:
			_ = d.Set("data_store", flattenCloudGcpDataStoreIntegration(t))
		case *cloud.CloudGcpFirebasedatabaseIntegration:
			_ = d.Set("fire_base_database", flattenCloudGcpFireBaseDatabaseIntegration(t))
		case *cloud.CloudGcpFirebasehostingIntegration:
			_ = d.Set("fire_base_hosting", flattenCloudGcpFireBaseHostingIntegration(t))
		case *cloud.CloudGcpFirebasestorageIntegration:
			_ = d.Set("fire_base_storage", flattenCloudGcpFireBaseStorageIntegration(t))
		case *cloud.CloudGcpFirestoreIntegration:
			_ = d.Set("fire_store", flattenCloudGcpFireStoreIntegration(t))
		case *cloud.CloudGcpFunctionsIntegration:
			_ = d.Set("functions", flattenCloudGcpFunctionsIntegration(t))
		case *cloud.CloudGcpInterconnectIntegration:
			_ = d.Set("interconnect", flattenCloudGcpInterconnectIntegration(t))
		case *cloud.CloudGcpKubernetesIntegration:
			_ = d.Set("kubernetes", flattenCloudGcpKubernetesIntegration(t))
		case *cloud.CloudGcpLoadbalancingIntegration:
			_ = d.Set("load_balancing", flattenCloudGcpLoadBalancingIntegration(t))
		case *cloud.CloudGcpMemcacheIntegration:
			_ = d.Set("mem_cache", flattenCloudGcpMemCacheIntegration(t))
		case *cloud.CloudGcpPubsubIntegration:
			_ = d.Set("pub_sub", flattenCloudGcpPubSubIntegration(t))
		case *cloud.CloudGcpRedisIntegration:
			_ = d.Set("redis", flattenCloudGcpRedisIntegration(t))
		case *cloud.CloudGcpRouterIntegration:
			_ = d.Set("router", flattenCloudGcpRouterIntegration(t))
		case *cloud.CloudGcpRunIntegration:
			_ = d.Set("run", flattenCloudGcpRunIntegration(t))
		case *cloud.CloudGcpSpannerIntegration:
			_ = d.Set("spanner", flattenCloudGcpSpannerIntegration(t))
		case *cloud.CloudGcpSqlIntegration:
			_ = d.Set("sql", flattenCloudGcpSQLIntegration(t))
		case *cloud.CloudGcpStorageIntegration:
			_ = d.Set("storage", flattenCloudGcpStorageIntegration(t))
		case *cloud.CloudGcpVmsIntegration:
			_ = d.Set("virtual_machines", flattenCloudGcpVirtualMachineIntegration(t))
		case *cloud.CloudGcpVpcaccessIntegration:
			_ = d.Set("vpc_access", flattenCloudGcpVpcAccessIntegration(t))
		}
	}
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpAppEngineIntegration(in *cloud.CloudGcpAppengineIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpBigQueryIntegration(in *cloud.CloudGcpBigqueryIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpBigTableIntegration(in *cloud.CloudGcpBigtableIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpComposerIntegration(in *cloud.CloudGcpComposerIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpDataFlowIntegration(in *cloud.CloudGcpDataflowIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpDataProcIntegration(in *cloud.CloudGcpDataprocIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpDataStoreIntegration(in *cloud.CloudGcpDatastoreIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpFireBaseDatabaseIntegration(in *cloud.CloudGcpFirebasedatabaseIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpFireBaseHostingIntegration(in *cloud.CloudGcpFirebasehostingIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpFireBaseStorageIntegration(in *cloud.CloudGcpFirebasestorageIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpFireStoreIntegration(in *cloud.CloudGcpFirestoreIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpFunctionsIntegration(in *cloud.CloudGcpFunctionsIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpInterconnectIntegration(in *cloud.CloudGcpInterconnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpKubernetesIntegration(in *cloud.CloudGcpKubernetesIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpLoadBalancingIntegration(in *cloud.CloudGcpLoadbalancingIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpMemCacheIntegration(in *cloud.CloudGcpMemcacheIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpPubSubIntegration(in *cloud.CloudGcpPubsubIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpRedisIntegration(in *cloud.CloudGcpRedisIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpRouterIntegration(in *cloud.CloudGcpRouterIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpRunIntegration(in *cloud.CloudGcpRunIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpSpannerIntegration(in *cloud.CloudGcpSpannerIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpSQLIntegration(in *cloud.CloudGcpSqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpStorageIntegration(in *cloud.CloudGcpStorageIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpVirtualMachineIntegration(in *cloud.CloudGcpVmsIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	flattened[0] = out
	return flattened
}

//flatten function to set(store) outputs from the terraform apply
func flattenCloudGcpVpcAccessIntegration(in *cloud.CloudGcpVpcaccessIntegration) []interface{} {
	flattened := make([]interface{}, 1)
	out := make(map[string]interface{})
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

//expand function to extract the inputs values from the schema for disabling the integration for particular services
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandCloudGcpDisableinputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudGcpDisableInput := cloud.CloudGcpDisableIntegrationsInput{}
	var linkedAccountID int
	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
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
