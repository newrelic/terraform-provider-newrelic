package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// gcpDmCheckLinkedAccountQuery is a minimal existence-check query that verifies
// the linked account still exists without requesting the integrations field.
// Requesting integrations causes "Abstract type 'Integration' must resolve to an
// Object type at runtime" errors on environments (e.g. staging) where GCP Dimensional Metrics-specific
// integration types are registered in the backend but not fully in the GraphQL schema.
const gcpDmCheckLinkedAccountQuery = `query($accountID: Int!, $id: Int!) {
	actor {
		account(id: $accountID) {
			cloud {
				linkedAccount(id: $id) {
					id
					nrAccountId
				}
			}
		}
	}
}`

// gcpDmCheckLinkedAccountResp is the response for gcpDmCheckLinkedAccountQuery.
type gcpDmCheckLinkedAccountResp struct {
	Actor struct {
		Account struct {
			Cloud struct {
				LinkedAccount *struct {
					ID          int `json:"id"`
					NrAccountId int `json:"nrAccountId"`
				} `json:"linkedAccount"`
			} `json:"cloud"`
		} `json:"account"`
	} `json:"actor"`
}

// gcpDmFilterDisableErrors filters out benign errors from a disable mutation
// response. Both ERR_INVALID_DATA and ERR_OBJECT_NOT_FOUND indicate the service
// was never enabled (slug unsupported on this environment, or integration simply
// doesn't exist), so there is nothing to disable — safe to skip.
// All other error types are aggregated and returned as a single error.
func gcpDmFilterDisableErrors(errors []cloud.CloudIntegrationMutationError) error {
	var fatal []string
	for _, e := range errors {
		if e.Type == "ERR_INVALID_DATA" || e.Type == "ERR_OBJECT_NOT_FOUND" {
			continue // service never enabled — nothing to disable, skip
		}
		fatal = append(fatal, e.Type+": "+e.Message)
	}
	if len(fatal) > 0 {
		return fmt.Errorf("cloudDisableIntegration errors: %s", strings.Join(fatal, "; "))
	}
	return nil
}

// ─── Resource definition ──────────────────────────────────────────────────────

func resourceNewrelicCloudGcpDmIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewrelicCloudGcpDmIntegrationsCreate,
		ReadContext:   resourceNewrelicCloudGcpDmIntegrationsRead,
		UpdateContext: resourceNewrelicCloudGcpDmIntegrationsUpdate,
		DeleteContext: resourceNewrelicCloudGcpDmIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: generateGcpDmIntegrationSchema(),
	}
}

func generateGcpDmIntegrationSchema() map[string]*schema.Schema {
	baseSchema := cloudGcpDmIntegrationSchemaBase()
	bigQuerySchema := cloudGcpDmMergeSchema(baseSchema, map[string]*schema.Schema{
		"fetch_tags": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Fetch resource tags for this integration.",
		},
		"fetch_table_metrics": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Fetch metrics for each table in BigQuery.",
		},
	})
	fetchTagsSchema := cloudGcpDmMergeSchema(baseSchema, map[string]*schema.Schema{
		"fetch_tags": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Fetch resource tags for this integration.",
		},
	})

	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The New Relic account ID.",
		},
		"linked_account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the GCP Dimensional Metrics linked account (from newrelic_cloud_gcp_dm_link_account).",
		},
		// ── Existing 27 services ──
		"ai_platform":       serviceBlock("GCP Vertex AI / AI Platform.", baseSchema),
		"alloy_db":          serviceBlock("GCP AlloyDB.", baseSchema),
		"app_engine":        serviceBlock("GCP App Engine.", baseSchema),
		"big_query":         serviceBlock("GCP BigQuery.", bigQuerySchema),
		"big_table":         serviceBlock("GCP Bigtable.", baseSchema),
		"composer":          serviceBlock("GCP Cloud Composer.", baseSchema),
		"data_flow":         serviceBlock("GCP Cloud Dataflow.", baseSchema),
		"data_proc":         serviceBlock("GCP Cloud Dataproc.", baseSchema),
		"data_store":        serviceBlock("GCP Cloud Datastore.", baseSchema),
		"firebase_database": serviceBlock("GCP Firebase Realtime Database.", baseSchema),
		"firebase_hosting":  serviceBlock("GCP Firebase Hosting.", baseSchema),
		"firebase_storage":  serviceBlock("GCP Firebase Storage.", baseSchema),
		"firestore":         serviceBlock("GCP Firestore.", baseSchema),
		"functions":         serviceBlock("GCP Cloud Functions.", baseSchema),
		"interconnect":      serviceBlock("GCP Cloud Interconnect.", baseSchema),
		"kubernetes":        serviceBlock("GCP Google Kubernetes Engine (GKE).", baseSchema),
		"load_balancing":    serviceBlock("GCP Cloud Load Balancing.", baseSchema),
		"mem_cache":         serviceBlock("GCP Memcache.", baseSchema),
		"pub_sub":           serviceBlock("GCP Cloud Pub/Sub.", fetchTagsSchema),
		"redis":             serviceBlock("GCP Memorystore for Redis (legacy).", baseSchema),
		"router":            serviceBlock("GCP Cloud Router.", baseSchema),
		"run":               serviceBlock("GCP Cloud Run.", baseSchema),
		"spanner":           serviceBlock("GCP Cloud Spanner.", fetchTagsSchema),
		"sql":               serviceBlock("GCP Cloud SQL.", baseSchema),
		"storage":           serviceBlock("GCP Cloud Storage.", fetchTagsSchema),
		"virtual_machines":  serviceBlock("GCP Compute Engine VMs.", baseSchema),
		"vpc_access":        serviceBlock("GCP Serverless VPC Access.", baseSchema),
		// ── New GCP Dimensional Metrics services ──
		"api_gateway":          serviceBlock("GCP API Gateway (Dimensional Metrics only).", baseSchema),
		"firebase_auth":        serviceBlock("Firebase Authentication (Dimensional Metrics only).", baseSchema),
		"firebase_vertex_ai":   serviceBlock("Firebase Vertex AI (Dimensional Metrics only; no entity synthesis).", baseSchema),
		"istio":                serviceBlock("GCP Istio Service Mesh (Dimensional Metrics only; no entity synthesis).", baseSchema),
		"managed_kafka":        serviceBlock("GCP Managed Service for Apache Kafka (Dimensional Metrics only).", baseSchema),
		"memory_store":         serviceBlock("GCP Memorystore for Redis/Memcached (Dimensional Metrics only).", baseSchema),
		"firebase_app_hosting": serviceBlock("Firebase App Hosting (Dimensional Metrics only; no entity synthesis).", baseSchema),
	}
}

// serviceBlock returns a TypeList schema.Schema with MaxItems:1 for a single integration block.
func serviceBlock(description string, elem map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: description,
		Elem:        &schema.Resource{Schema: elem},
	}
}

// cloudGcpDmIntegrationSchemaBase is the minimal schema shared by all GCP Dimensional Metrics service blocks.
func cloudGcpDmIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds.",
		},
	}
}

// cloudGcpDmMergeSchema merges two schema maps into a new map (non-destructive).
func cloudGcpDmMergeSchema(base, extra map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema, len(base)+len(extra))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range extra {
		result[k] = v
	}
	return result
}

// ─── CRUD functions ───────────────────────────────────────────────────────────

func resourceNewrelicCloudGcpDmIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID := d.Get("linked_account_id").(int)

	gcpInput, _ := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	configPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloud.CloudIntegrationsInput{Gcp: gcpInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration failed: %w", err))
	}
	if len(configPayload.Errors) > 0 {
		msgs := make([]string, 0, len(configPayload.Errors))
		for _, e := range configPayload.Errors {
			msgs = append(msgs, e.Type+": "+e.Message)
		}
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration errors: %s", strings.Join(msgs, "; ")))
	}

	d.SetId(strconv.Itoa(linkedAccountID))
	_ = d.Set("account_id", accountID)

	return resourceNewrelicCloudGcpDmIntegrationsRead(ctx, d, meta)
}

func resourceNewrelicCloudGcpDmIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Existence-check only — do not request integrations.
	// Fetching the integrations field causes "Abstract type 'Integration' must resolve
	// to an Object type" errors on environments where GCP Dimensional Metrics integration types are not
	// fully registered in the GraphQL schema (e.g. staging). The linked_account_id
	// and account_id are already in state; this Read simply confirms the account exists.
	var checkResp gcpDmCheckLinkedAccountResp
	vars := map[string]interface{}{
		"accountID": accountID,
		"id":        linkedAccountID,
	}
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, gcpDmCheckLinkedAccountQuery, vars, &checkResp); err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("linkedAccount existence check failed: %w", err))
	}

	if checkResp.Actor.Account.Cloud.LinkedAccount == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", checkResp.Actor.Account.Cloud.LinkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccountID)

	return nil
}

func resourceNewrelicCloudGcpDmIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gcpInput, gcpDisable := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	// Disable removed integrations first
	disablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, cloud.CloudDisableIntegrationsInput{Gcp: gcpDisable})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudDisableIntegration failed: %w", err))
	}
	if err := gcpDmFilterDisableErrors(disablePayload.Errors); err != nil {
		return diag.FromErr(err)
	}

	// Enable/update present integrations
	configPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloud.CloudIntegrationsInput{Gcp: gcpInput})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration failed: %w", err))
	}
	if len(configPayload.Errors) > 0 {
		msgs := make([]string, 0, len(configPayload.Errors))
		for _, e := range configPayload.Errors {
			msgs = append(msgs, e.Type+": "+e.Message)
		}
		return diag.FromErr(fmt.Errorf("cloudConfigureIntegration errors: %s", strings.Join(msgs, "; ")))
	}

	return resourceNewrelicCloudGcpDmIntegrationsRead(ctx, d, meta)
}

func resourceNewrelicCloudGcpDmIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, gcpDisable := expandCloudGcpDmIntegrationsInput(d, linkedAccountID)

	disablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, cloud.CloudDisableIntegrationsInput{Gcp: gcpDisable})
	if err != nil {
		return diag.FromErr(fmt.Errorf("cloudDisableIntegration failed: %w", err))
	}
	if err := gcpDmFilterDisableErrors(disablePayload.Errors); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// ─── Expand function ──────────────────────────────────────────────────────────

// expandCloudGcpDmIntegrationsInput builds configure and disable inputs for all 34 GCP services.
// Present blocks go to configure; absent blocks go to disable.
// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func expandCloudGcpDmIntegrationsInput(d *schema.ResourceData, linkedAccountID int) (cloud.CloudGcpIntegrationsInput, cloud.CloudGcpDisableIntegrationsInput) {
	gcpInput := cloud.CloudGcpIntegrationsInput{}
	gcpDisable := cloud.CloudGcpDisableIntegrationsInput{}
	dis := cloud.CloudDisableAccountIntegrationInput{LinkedAccountId: linkedAccountID}

	present := func(key string) bool {
		v, ok := d.GetOk(key)
		if !ok {
			return false
		}
		l, ok := v.([]interface{})
		return ok && len(l) > 0
	}

	getInt := func(key string) int {
		if v := d.Get(key); v != nil {
			return v.(int)
		}
		return 0
	}

	getBool := func(key string) bool {
		if v := d.Get(key); v != nil {
			return v.(bool)
		}
		return false
	}

	// ── Existing 27 services ─────────────────────────────────────────────────

	if present("ai_platform") {
		gcpInput.GcpAiplatform = []cloud.CloudGcpAiplatformIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("ai_platform.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpAiplatform = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("alloy_db") {
		gcpInput.GcpAlloydb = []cloud.CloudGcpAlloydbIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("alloy_db.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("app_engine") {
		gcpInput.GcpAppengine = []cloud.CloudGcpAppengineIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("app_engine.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("big_query") {
		gcpInput.GcpBigquery = []cloud.CloudGcpBigqueryIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("big_query.0.metrics_polling_interval"),
			FetchTags:              getBool("big_query.0.fetch_tags"),
			FetchTableMetrics:      getBool("big_query.0.fetch_table_metrics"),
		}}
	} else {
		gcpDisable.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("big_table") {
		gcpInput.GcpBigtable = []cloud.CloudGcpBigtableIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("big_table.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("composer") {
		gcpInput.GcpComposer = []cloud.CloudGcpComposerIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("composer.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpComposer = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("data_flow") {
		gcpInput.GcpDataflow = []cloud.CloudGcpDataflowIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("data_flow.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpDataflow = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("data_proc") {
		gcpInput.GcpDataproc = []cloud.CloudGcpDataprocIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("data_proc.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpDataproc = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("data_store") {
		gcpInput.GcpDatastore = []cloud.CloudGcpDatastoreIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("data_store.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpDatastore = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_database") {
		gcpInput.GcpFirebasedatabase = []cloud.CloudGcpFirebasedatabaseIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_database.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebasedatabase = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_hosting") {
		gcpInput.GcpFirebasehosting = []cloud.CloudGcpFirebasehostingIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_hosting.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebasehosting = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_storage") {
		gcpInput.GcpFirebasestorage = []cloud.CloudGcpFirebasestorageIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_storage.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebasestorage = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firestore") {
		gcpInput.GcpFirestore = []cloud.CloudGcpFirestoreIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firestore.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirestore = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("functions") {
		gcpInput.GcpFunctions = []cloud.CloudGcpFunctionsIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("functions.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFunctions = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("interconnect") {
		gcpInput.GcpInterconnect = []cloud.CloudGcpInterconnectIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("interconnect.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpInterconnect = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("kubernetes") {
		gcpInput.GcpKubernetes = []cloud.CloudGcpKubernetesIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("kubernetes.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpKubernetes = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("load_balancing") {
		gcpInput.GcpLoadbalancing = []cloud.CloudGcpLoadbalancingIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("load_balancing.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpLoadbalancing = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("mem_cache") {
		gcpInput.GcpMemcache = []cloud.CloudGcpMemcacheIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("mem_cache.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpMemcache = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("pub_sub") {
		gcpInput.GcpPubsub = []cloud.CloudGcpPubsubIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("pub_sub.0.metrics_polling_interval"),
			FetchTags:              getBool("pub_sub.0.fetch_tags"),
		}}
	} else {
		gcpDisable.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("redis") {
		gcpInput.GcpRedis = []cloud.CloudGcpRedisIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("redis.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpRedis = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("router") {
		gcpInput.GcpRouter = []cloud.CloudGcpRouterIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("router.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpRouter = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("run") {
		gcpInput.GcpRun = []cloud.CloudGcpRunIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("run.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpRun = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("spanner") {
		gcpInput.GcpSpanner = []cloud.CloudGcpSpannerIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("spanner.0.metrics_polling_interval"),
			FetchTags:              getBool("spanner.0.fetch_tags"),
		}}
	} else {
		gcpDisable.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("sql") {
		gcpInput.GcpSql = []cloud.CloudGcpSqlIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("sql.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpSql = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("storage") {
		gcpInput.GcpStorage = []cloud.CloudGcpStorageIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("storage.0.metrics_polling_interval"),
			FetchTags:              getBool("storage.0.fetch_tags"),
		}}
	} else {
		gcpDisable.GcpStorage = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("virtual_machines") {
		gcpInput.GcpVms = []cloud.CloudGcpVmsIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("virtual_machines.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpVms = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("vpc_access") {
		gcpInput.GcpVpcaccess = []cloud.CloudGcpVpcaccessIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("vpc_access.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpVpcaccess = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	// ── New GCP Dimensional Metrics services ─────────────────────────────────────────────────

	if present("api_gateway") {
		gcpInput.GcpApiGateway = []cloud.CloudGcpApiGatewayIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("api_gateway.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpApiGateway = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_auth") {
		gcpInput.GcpFirebaseAuth = []cloud.CloudGcpFirebaseAuthIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_auth.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebaseAuth = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_vertex_ai") {
		gcpInput.GcpFirebaseVertexAi = []cloud.CloudGcpFirebaseVertexAiIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_vertex_ai.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebaseVertexAi = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("istio") {
		gcpInput.GcpIstio = []cloud.CloudGcpIstioIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("istio.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpIstio = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("managed_kafka") {
		gcpInput.GcpManagedKafka = []cloud.CloudGcpManagedKafkaIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("managed_kafka.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpManagedKafka = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("memory_store") {
		gcpInput.GcpMemoryStore = []cloud.CloudGcpMemoryStoreIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("memory_store.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpMemoryStore = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	if present("firebase_app_hosting") {
		gcpInput.GcpFirebaseAppHosting = []cloud.CloudGcpFirebaseAppHostingIntegrationInput{{
			LinkedAccountId:        linkedAccountID,
			MetricsPollingInterval: getInt("firebase_app_hosting.0.metrics_polling_interval"),
		}}
	} else {
		gcpDisable.GcpFirebaseAppHosting = []cloud.CloudDisableAccountIntegrationInput{dis}
	}

	return gcpInput, gcpDisable
}

