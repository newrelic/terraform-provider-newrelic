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

// ─── Service table ──────────────────────────────────────────────────────────
//
// Every GCP Dimensional Metrics service is described by a single table entry. The
// schema, the configure input, and the disable input are all derived from this
// table, so adding a service is a one-line change rather than edits to three
// parallel switch/if-else blocks.

// gcpDmServiceValues carries the per-service values read from the resource schema.
// GCP Dimensional Metrics is metrics-only, so metrics_polling_interval is the only
// per-service knob — the backend rejects inventory-era params such as fetch_tags.
type gcpDmServiceValues struct {
	LinkedAccountID        int
	MetricsPollingInterval int
}

// gcpDmService describes one GCP service block.
type gcpDmService struct {
	key         string
	description string
	// configure sets this service's typed input slice on the configure payload.
	configure func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues)
	// disable sets this service's disable slice on the disable payload.
	disable func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput)
}

// genericConfigure builds a configure closure for a DM-only service backed by the
// consolidated CloudGcpGenericIntegrationInput type (only linked account + polling).
func genericConfigure(assign func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput)) func(*cloud.CloudGcpIntegrationsInput, gcpDmServiceValues) {
	return func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
		assign(g, []cloud.CloudGcpGenericIntegrationInput{{
			LinkedAccountId:        v.LinkedAccountID,
			MetricsPollingInterval: v.MetricsPollingInterval,
		}})
	}
}

// gcpDmServices returns the full catalog of supported GCP Dimensional Metrics services:
// the 27 services shared with the legacy resource plus the 7 DM-only services.
func gcpDmServices() []gcpDmService {
	return []gcpDmService{
		// ── Shared 27 services (per-service typed inputs) ────────────────────
		{key: "ai_platform", description: "GCP Vertex AI / AI Platform.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpAiplatform = []cloud.CloudGcpAiplatformIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpAiplatform = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "alloy_db", description: "GCP AlloyDB.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpAlloydb = []cloud.CloudGcpAlloydbIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "app_engine", description: "GCP App Engine.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpAppengine = []cloud.CloudGcpAppengineIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "big_query", description: "GCP BigQuery.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpBigquery = []cloud.CloudGcpBigqueryIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "big_table", description: "GCP Bigtable.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpBigtable = []cloud.CloudGcpBigtableIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "composer", description: "GCP Cloud Composer.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpComposer = []cloud.CloudGcpComposerIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpComposer = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "data_flow", description: "GCP Cloud Dataflow.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpDataflow = []cloud.CloudGcpDataflowIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpDataflow = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "data_proc", description: "GCP Cloud Dataproc.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpDataproc = []cloud.CloudGcpDataprocIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpDataproc = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "data_store", description: "GCP Cloud Datastore.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpDatastore = []cloud.CloudGcpDatastoreIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpDatastore = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_database", description: "GCP Firebase Realtime Database.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpFirebasedatabase = []cloud.CloudGcpFirebasedatabaseIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebasedatabase = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_hosting", description: "GCP Firebase Hosting.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpFirebasehosting = []cloud.CloudGcpFirebasehostingIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebasehosting = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_storage", description: "GCP Firebase Storage.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpFirebasestorage = []cloud.CloudGcpFirebasestorageIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebasestorage = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firestore", description: "GCP Firestore.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpFirestore = []cloud.CloudGcpFirestoreIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirestore = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "functions", description: "GCP Cloud Functions.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpFunctions = []cloud.CloudGcpFunctionsIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFunctions = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "interconnect", description: "GCP Cloud Interconnect.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpInterconnect = []cloud.CloudGcpInterconnectIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpInterconnect = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "kubernetes", description: "GCP Google Kubernetes Engine (GKE).",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpKubernetes = []cloud.CloudGcpKubernetesIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpKubernetes = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "load_balancing", description: "GCP Cloud Load Balancing.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpLoadbalancing = []cloud.CloudGcpLoadbalancingIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpLoadbalancing = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "mem_cache", description: "GCP Memcache.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpMemcache = []cloud.CloudGcpMemcacheIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpMemcache = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "pub_sub", description: "GCP Cloud Pub/Sub.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpPubsub = []cloud.CloudGcpPubsubIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "redis", description: "GCP Memorystore for Redis (legacy).",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpRedis = []cloud.CloudGcpRedisIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpRedis = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "router", description: "GCP Cloud Router.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpRouter = []cloud.CloudGcpRouterIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpRouter = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "run", description: "GCP Cloud Run.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpRun = []cloud.CloudGcpRunIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpRun = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "spanner", description: "GCP Cloud Spanner.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpSpanner = []cloud.CloudGcpSpannerIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "sql", description: "GCP Cloud SQL.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpSql = []cloud.CloudGcpSqlIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpSql = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "storage", description: "GCP Cloud Storage.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpStorage = []cloud.CloudGcpStorageIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpStorage = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "virtual_machines", description: "GCP Compute Engine VMs.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpVms = []cloud.CloudGcpVmsIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpVms = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "vpc_access", description: "GCP Serverless VPC Access.",
			configure: func(g *cloud.CloudGcpIntegrationsInput, v gcpDmServiceValues) {
				g.GcpVpcaccess = []cloud.CloudGcpVpcaccessIntegrationInput{{LinkedAccountId: v.LinkedAccountID, MetricsPollingInterval: v.MetricsPollingInterval}}
			},
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpVpcaccess = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},

		// ── DM-only 7 services (consolidated CloudGcpGenericIntegrationInput) ──
		{key: "api_gateway", description: "GCP API Gateway (Dimensional Metrics only).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpApiGateway = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpApiGateway = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_auth", description: "Firebase Authentication (Dimensional Metrics only).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpFirebaseAuth = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebaseAuth = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_vertex_ai", description: "Firebase Vertex AI (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpFirebaseVertexAi = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebaseVertexAi = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "istio", description: "GCP Istio Service Mesh (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) { g.GcpIstio = in }),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpIstio = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "managed_kafka", description: "GCP Managed Service for Apache Kafka (Dimensional Metrics only).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpManagedKafka = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpManagedKafka = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "memory_store", description: "GCP Memorystore for Redis/Memcached (Dimensional Metrics only).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpMemoryStore = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpMemoryStore = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
		{key: "firebase_app_hosting", description: "Firebase App Hosting (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(g *cloud.CloudGcpIntegrationsInput, in []cloud.CloudGcpGenericIntegrationInput) {
				g.GcpFirebaseAppHosting = in
			}),
			disable: func(g *cloud.CloudGcpDisableIntegrationsInput, dis cloud.CloudDisableAccountIntegrationInput) {
				g.GcpFirebaseAppHosting = []cloud.CloudDisableAccountIntegrationInput{dis}
			}},
	}
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
	base := cloudGcpDmIntegrationSchemaBase()

	s := map[string]*schema.Schema{
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
			Description: "The ID of the GCP Dimensional Metrics linked account (from newrelic_cloud_gcp_link_account with audience + service_account_email set).",
		},
	}

	for _, svc := range gcpDmServices() {
		s[svc.key] = serviceBlock(svc.description, base)
	}

	return s
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

// expandCloudGcpDmIntegrationsInput builds configure and disable inputs for all GCP
// services from the service table. Present blocks go to configure; absent blocks go
// to disable.
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

	for _, svc := range gcpDmServices() {
		if !present(svc.key) {
			svc.disable(&gcpDisable, dis)
			continue
		}
		svc.configure(&gcpInput, gcpDmServiceValues{
			LinkedAccountID:        linkedAccountID,
			MetricsPollingInterval: getInt(svc.key + ".0.metrics_polling_interval"),
		})
	}

	return gcpInput, gcpDisable
}
