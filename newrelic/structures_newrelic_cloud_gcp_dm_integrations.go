package newrelic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// ─── Service table ──────────────────────────────────────────────────────────
//
// Every GCP Dimensional Metrics service is described by a single table entry.
// The schema, the configure input, the disable input, AND the read/flatten are
// all derived from this one table, so adding a service is a one-line change.
//
// On a GCP Dimensional Metrics (gcp_v2) linked account, the backend returns every
// integration as a CloudGcpGenericIntegration identified by its service `slug`
// (e.g. "gcp_bigquery"); the `slug` field below is how the read maps a returned
// integration back to its Terraform block.

// gcpDmServiceValues carries the per-service values read from the resource schema.
// GCP Dimensional Metrics is metrics-only, so metrics_polling_interval is the only
// per-service knob — the backend rejects inventory-era params such as fetch_tags.
type gcpDmServiceValues struct {
	LinkedAccountID        int
	MetricsPollingInterval int
}

// gcpDmService describes one GCP service block.
type gcpDmService struct {
	key         string // Terraform schema block key, e.g. "big_query"
	slug        string // backend CloudService.Slug on gcp_v2 accounts, e.g. "gcp_bigquery"
	description string
	// configure sets this service's typed input slice on the configure payload.
	configure func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues)
	// disable sets this service's disable slice on the disable payload.
	disable func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput)
}

// genericConfigure builds a configure closure for a DM-only service backed by the
// consolidated CloudGcpGenericIntegrationInput type (only linked account + polling).
func genericConfigure(assign func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput)) func(*cloud.CloudGcpIntegrationsInput, gcpDmServiceValues) {
	return func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
		assign(in, []cloud.CloudGcpGenericIntegrationInput{{
			LinkedAccountId:        values.LinkedAccountID,
			MetricsPollingInterval: values.MetricsPollingInterval,
		}})
	}
}

// gcpDmServices returns the full catalog of supported GCP Dimensional Metrics services:
// the 27 services shared with the legacy resource plus the 7 DM-only services.
func gcpDmServices() []gcpDmService {
	return []gcpDmService{
		// ── Shared 27 services (per-service typed configure inputs) ──────────
		{key: "ai_platform", slug: "gcp_aiplatform", description: "GCP Vertex AI / AI Platform.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpAiplatform = []cloud.CloudGcpAiplatformIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpAiplatform = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "alloy_db", slug: "gcp_alloydb", description: "GCP AlloyDB.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpAlloydb = []cloud.CloudGcpAlloydbIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "app_engine", slug: "gcp_appengine", description: "GCP App Engine.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpAppengine = []cloud.CloudGcpAppengineIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "big_query", slug: "gcp_bigquery", description: "GCP BigQuery.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpBigquery = []cloud.CloudGcpBigqueryIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "big_table", slug: "gcp_bigtable", description: "GCP Bigtable.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpBigtable = []cloud.CloudGcpBigtableIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "composer", slug: "gcp_composer", description: "GCP Cloud Composer.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpComposer = []cloud.CloudGcpComposerIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpComposer = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "data_flow", slug: "gcp_dataflow", description: "GCP Cloud Dataflow.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpDataflow = []cloud.CloudGcpDataflowIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpDataflow = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "data_proc", slug: "gcp_dataproc", description: "GCP Cloud Dataproc.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpDataproc = []cloud.CloudGcpDataprocIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpDataproc = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "data_store", slug: "gcp_datastore", description: "GCP Cloud Datastore.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpDatastore = []cloud.CloudGcpDatastoreIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpDatastore = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_database", slug: "gcp_firebasedatabase", description: "GCP Firebase Realtime Database.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpFirebasedatabase = []cloud.CloudGcpFirebasedatabaseIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebasedatabase = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_hosting", slug: "gcp_firebasehosting", description: "GCP Firebase Hosting.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpFirebasehosting = []cloud.CloudGcpFirebasehostingIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebasehosting = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_storage", slug: "gcp_firebasestorage", description: "GCP Firebase Storage.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpFirebasestorage = []cloud.CloudGcpFirebasestorageIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebasestorage = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firestore", slug: "gcp_firestore", description: "GCP Firestore.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpFirestore = []cloud.CloudGcpFirestoreIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirestore = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "functions", slug: "gcp_functions", description: "GCP Cloud Functions.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpFunctions = []cloud.CloudGcpFunctionsIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFunctions = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "interconnect", slug: "gcp_interconnect", description: "GCP Cloud Interconnect.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpInterconnect = []cloud.CloudGcpInterconnectIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpInterconnect = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "kubernetes", slug: "gcp_kubernetes", description: "GCP Google Kubernetes Engine (GKE).",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpKubernetes = []cloud.CloudGcpKubernetesIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpKubernetes = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "load_balancing", slug: "gcp_loadbalancing", description: "GCP Cloud Load Balancing.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpLoadbalancing = []cloud.CloudGcpLoadbalancingIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpLoadbalancing = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "mem_cache", slug: "gcp_memcache", description: "GCP Memcache.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpMemcache = []cloud.CloudGcpMemcacheIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpMemcache = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "pub_sub", slug: "gcp_pubsub", description: "GCP Cloud Pub/Sub.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpPubsub = []cloud.CloudGcpPubsubIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "redis", slug: "gcp_redis", description: "GCP Memorystore for Redis (legacy).",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpRedis = []cloud.CloudGcpRedisIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpRedis = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "router", slug: "gcp_router", description: "GCP Cloud Router.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpRouter = []cloud.CloudGcpRouterIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpRouter = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "run", slug: "gcp_run", description: "GCP Cloud Run.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpRun = []cloud.CloudGcpRunIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpRun = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "spanner", slug: "gcp_spanner", description: "GCP Cloud Spanner.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpSpanner = []cloud.CloudGcpSpannerIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "sql", slug: "gcp_sql", description: "GCP Cloud SQL.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpSql = []cloud.CloudGcpSqlIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpSql = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "storage", slug: "gcp_storage", description: "GCP Cloud Storage.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpStorage = []cloud.CloudGcpStorageIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpStorage = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "virtual_machines", slug: "gcp_vms", description: "GCP Compute Engine VMs.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpVms = []cloud.CloudGcpVmsIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpVms = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "vpc_access", slug: "gcp_vpcaccess", description: "GCP Serverless VPC Access.",
			configure: func(in *cloud.CloudGcpIntegrationsInput, values gcpDmServiceValues) {
				in.GcpVpcaccess = []cloud.CloudGcpVpcaccessIntegrationInput{{LinkedAccountId: values.LinkedAccountID, MetricsPollingInterval: values.MetricsPollingInterval}}
			},
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpVpcaccess = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},

		// ── DM-only 7 services (consolidated CloudGcpGenericIntegrationInput) ──
		{key: "api_gateway", slug: "gcp_api_gateway", description: "GCP API Gateway (Dimensional Metrics only).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpApiGateway = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpApiGateway = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_auth", slug: "gcp_firebase_auth", description: "Firebase Authentication (Dimensional Metrics only).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpFirebaseAuth = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebaseAuth = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_vertex_ai", slug: "gcp_firebase_vertex_ai", description: "Firebase Vertex AI (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpFirebaseVertexAi = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebaseVertexAi = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "istio", slug: "gcp_istio", description: "GCP Istio Service Mesh (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpIstio = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpIstio = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "managed_kafka", slug: "gcp_managed_kafka", description: "GCP Managed Service for Apache Kafka (Dimensional Metrics only).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpManagedKafka = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpManagedKafka = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "memory_store", slug: "gcp_memory_store", description: "GCP Memorystore for Redis/Memcached (Dimensional Metrics only).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpMemoryStore = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpMemoryStore = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
		{key: "firebase_app_hosting", slug: "gcp_firebase_app_hosting", description: "Firebase App Hosting (Dimensional Metrics only; no entity synthesis).",
			configure: genericConfigure(func(in *cloud.CloudGcpIntegrationsInput, generic []cloud.CloudGcpGenericIntegrationInput) {
				in.GcpFirebaseAppHosting = generic
			}),
			disable: func(in *cloud.CloudGcpDisableIntegrationsInput, disableInput cloud.CloudDisableAccountIntegrationInput) {
				in.GcpFirebaseAppHosting = []cloud.CloudDisableAccountIntegrationInput{disableInput}
			}},
	}
}

// gcpDmSlugToKey maps each backend service slug to its Terraform schema block key,
// derived from the single service table above. Used by the read/flatten to place a
// returned CloudGcpGenericIntegration into the right block.
func gcpDmSlugToKey() map[string]string {
	m := make(map[string]string, len(gcpDmServices()))
	for _, service := range gcpDmServices() {
		m[service.slug] = service.key
	}
	return m
}

// ─── Schema ─────────────────────────────────────────────────────────────────

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
			Description: "The ID of the GCP Dimensional Metrics linked account (from newrelic_cloud_gcp_link_account with use_workload_identity_federation = true).",
		},
	}

	for _, service := range gcpDmServices() {
		s[service.key] = serviceBlock(service.description, base)
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

// cloudGcpDmIntegrationSchemaBase is the minimal schema shared by all GCP Dimensional
// Metrics service blocks. metrics_polling_interval is Computed so that omitting it
// accepts whatever interval the backend reports (no spurious drift on read).
func cloudGcpDmIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The data polling interval in seconds.",
		},
	}
}

// ─── Expand / flatten ───────────────────────────────────────────────────────

// gcpDmFilterDisableErrors filters out benign errors from a disable mutation
// response. Both ERR_INVALID_DATA and ERR_OBJECT_NOT_FOUND indicate the service
// was never enabled (slug unsupported on this environment, or integration simply
// doesn't exist), so there is nothing to disable — safe to skip. All other error
// types are aggregated and returned as a single error.
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

// expandCloudGcpDmIntegrationsInput builds the configure and disable payloads from
// the resource config, driven by the service table. A service block present in the
// config is added to the configure payload (enable/update); an absent block is added
// to the disable payload (so removing a block from config disables that service).
func expandCloudGcpDmIntegrationsInput(d *schema.ResourceData, linkedAccountID int) (cloud.CloudGcpIntegrationsInput, cloud.CloudGcpDisableIntegrationsInput) {
	configureInput := cloud.CloudGcpIntegrationsInput{}
	disableInput := cloud.CloudGcpDisableIntegrationsInput{}
	disableForAccount := cloud.CloudDisableAccountIntegrationInput{LinkedAccountId: linkedAccountID}

	blockPresent := func(key string) bool {
		v, ok := d.GetOk(key)
		if !ok {
			return false
		}
		list, ok := v.([]interface{})
		return ok && len(list) > 0
	}

	getInt := func(key string) int {
		if v := d.Get(key); v != nil {
			return v.(int)
		}
		return 0
	}

	for _, service := range gcpDmServices() {
		if !blockPresent(service.key) {
			service.disable(&disableInput, disableForAccount) // absent → disable
			continue
		}
		service.configure(&configureInput, gcpDmServiceValues{ // present → enable/update
			LinkedAccountID:        linkedAccountID,
			MetricsPollingInterval: getInt(service.key + ".0.metrics_polling_interval"),
		})
	}

	return configureInput, disableInput
}

// flattenGcpDmIntegrations writes the backend state back into the resource so drift
// is detected. On a gcp_v2 linked account every integration is returned as a
// *cloud.CloudGcpGenericIntegration identified by its service slug; each is mapped
// back to its Terraform block. Blocks with no matching integration are cleared so an
// integration removed outside Terraform surfaces as drift.
func flattenGcpDmIntegrations(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	slugToKey := gcpDmSlugToKey()

	// Clear every service block first; only the ones present on the account are re-set.
	for _, service := range gcpDmServices() {
		_ = d.Set(service.key, []interface{}{})
	}

	for _, integration := range linkedAccount.Integrations {
		generic, ok := integration.(*cloud.CloudGcpGenericIntegration)
		if !ok {
			continue
		}
		key, ok := slugToKey[generic.Service.Slug]
		if !ok {
			continue
		}
		_ = d.Set(key, []interface{}{map[string]interface{}{
			"metrics_polling_interval": generic.MetricsPollingInterval,
		}})
	}
}
