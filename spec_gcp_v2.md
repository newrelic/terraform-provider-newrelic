# GCP v2 Terraform Provider Integration — Implementation Specification

> **Self-contained spec for a separate implementation thread.**
> All referenced file paths are absolute. All type names, field names, and method signatures are exact.

---

## 1. Overview

GCP Cloud Integrations v2 moves authentication from service account key files (GCP v1) to **Workload Identity Federation (WIF)** — short-lived credentials tied to an identity pool with no long-lived secret keys.

This spec adds three Terraform artefacts:

| Artefact | Type | Purpose |
|---|---|---|
| `newrelic_cloud_gcp_v2_link_account` | Resource | Link a GCP project to New Relic via WIF (2-step flow) |
| `newrelic_cloud_gcp_v2_integrations` | Resource | Configure polling integrations (27 existing + 7 new services) |
| `newrelic_cloud_gcp_v2_account` | Data Source | Look up a GCP v2 linked account ID by name |

**Existing resources are not modified.** `newrelic_cloud_gcp_link_account` and `newrelic_cloud_gcp_integrations` remain exactly as they are.

---

## 2. Background & Terminology

| Term | Meaning |
|---|---|
| **GCP v1** | Service account key auth; `provider_slug = "gcp"` in backend DB; sample-based metrics |
| **GCP v2** | WIF auth; `provider_slug = "gcp_v2"` in backend DB (internal only); dimensional metrics |
| **WIF** | Workload Identity Federation — customer provides a GCP-generated JSON credential |
| **authReferenceId** | UUID session key returned by `cloudAuthenticateIntegration`; 30-min TTL in Redis |
| **GcpGenericIntegration** | API type for the 7 new GCP v2 services (polling interval only, no service-specific params) |

`gcp_v2` is an **internal** backend concept. API callers never pass `"gcp_v2"` as a provider slug to public-facing queries. The provider slug is set transparently when `authReferenceId` is present in `linkAccount`.

---

## 3. Dependency: Changes Required in `newrelic-client-go`

> **Repo**: `github.com/newrelic/newrelic-client-go`
> **Current version in provider**: `v2.83.0`
> **Files to change**: `pkg/cloud/types.go`, `pkg/cloud/cloud_api.go`, `pkg/cloud/cloud.go`

After modifying client-go, update `go.mod` in the Terraform provider to point to the new version (or use a `replace` directive during development).

### 3.1 New Types — add to `pkg/cloud/types.go`

```go
// CloudGcpGenericIntegrationInput — input for new GCP v2 services
// (gcpApiGateway, gcpFirebaseAuth, gcpFirebaseVertexAi, gcpIstio,
//  gcpManagedKafka, gcpMemoryStore, gcpFirebaseAppHosting)
type CloudGcpGenericIntegrationInput struct {
    LinkedAccountId        int `json:"linkedAccountId"`
    MetricsPollingInterval int `json:"metricsPollingInterval,omitempty"`
}

// CloudGcpGenericIntegration — output type for new GCP v2 services
type CloudGcpGenericIntegration struct {
    CreatedAt              nrtime.EpochSeconds `json:"createdAt"`
    ID                     int                 `json:"id,omitempty"`
    LinkedAccount          CloudLinkedAccount  `json:"linkedAccount,omitempty"`
    MetricsPollingInterval int                 `json:"metricsPollingInterval,omitempty"`
    Name                   string              `json:"name,omitempty"`
    NrAccountId            int                 `json:"nrAccountId"`
    Service                CloudService        `json:"service,omitempty"`
    UpdatedAt              nrtime.EpochSeconds `json:"updatedAt"`
}

// CloudAuthenticateIntegrationInput — input for WIF step 1
type CloudAuthenticateIntegrationInput struct {
    NrAccountId  int    `json:"nrAccountId"`
    ProviderSlug string `json:"providerSlug"` // use "GCP"
    AuthType     string `json:"authType"`     // use "WIF"
    Payload      string `json:"payload"`      // WIF JSON credential as string
}

// CloudAuthenticateIntegrationPayload — response from WIF step 1
type CloudAuthenticateIntegrationPayload struct {
    AuthReferenceId string `json:"authReferenceId"`
}
```

### 3.2 Extend `CloudGcpLinkAccountInput` — modify in `pkg/cloud/types.go`

```go
// BEFORE
type CloudGcpLinkAccountInput struct {
    Name      string `json:"name"`
    ProjectId string `json:"projectId"`
}

// AFTER — add AuthReferenceId (omitempty keeps v1 calls unchanged)
type CloudGcpLinkAccountInput struct {
    Name            string `json:"name"`
    ProjectId       string `json:"projectId"`
    AuthReferenceId string `json:"authReferenceId,omitempty"`
}
```

### 3.3 Extend `CloudGcpIntegrationsInput` — modify in `pkg/cloud/types.go`

Append 7 new fields to the existing struct (do not change or reorder existing fields):

```go
// Append these 7 fields to CloudGcpIntegrationsInput:
GcpApiGateway         []CloudGcpGenericIntegrationInput `json:"gcpApiGateway,omitempty"`
GcpFirebaseAuth       []CloudGcpGenericIntegrationInput `json:"gcpFirebaseAuth,omitempty"`
GcpFirebaseVertexAi   []CloudGcpGenericIntegrationInput `json:"gcpFirebaseVertexAi,omitempty"`
GcpIstio              []CloudGcpGenericIntegrationInput `json:"gcpIstio,omitempty"`
GcpManagedKafka       []CloudGcpGenericIntegrationInput `json:"gcpManagedKafka,omitempty"`
GcpMemoryStore        []CloudGcpGenericIntegrationInput `json:"gcpMemoryStore,omitempty"`
GcpFirebaseAppHosting []CloudGcpGenericIntegrationInput `json:"gcpFirebaseAppHosting,omitempty"`
```

### 3.4 Extend `CloudGcpDisableIntegrationsInput` — modify in `pkg/cloud/types.go`

Append the same 7 fields (all use `CloudDisableAccountIntegrationInput`):

```go
// Append these 7 fields to CloudGcpDisableIntegrationsInput:
GcpApiGateway         []CloudDisableAccountIntegrationInput `json:"gcpApiGateway,omitempty"`
GcpFirebaseAuth       []CloudDisableAccountIntegrationInput `json:"gcpFirebaseAuth,omitempty"`
GcpFirebaseVertexAi   []CloudDisableAccountIntegrationInput `json:"gcpFirebaseVertexAi,omitempty"`
GcpIstio              []CloudDisableAccountIntegrationInput `json:"gcpIstio,omitempty"`
GcpManagedKafka       []CloudDisableAccountIntegrationInput `json:"gcpManagedKafka,omitempty"`
GcpMemoryStore        []CloudDisableAccountIntegrationInput `json:"gcpMemoryStore,omitempty"`
GcpFirebaseAppHosting []CloudDisableAccountIntegrationInput `json:"gcpFirebaseAppHosting,omitempty"`
```

### 3.5 New Method — add to `pkg/cloud/cloud.go`

```go
func (a *Cloud) CloudAuthenticateIntegrationWithContext(
    ctx context.Context,
    nrAccountId int,
    providerSlug string,
    authType string,
    payload string,
) (*CloudAuthenticateIntegrationPayload, error) {
    resp := &struct {
        CloudAuthenticateIntegration CloudAuthenticateIntegrationPayload `json:"cloudAuthenticateIntegration"`
    }{}
    vars := map[string]interface{}{
        "nrAccountId":  nrAccountId,
        "providerSlug": providerSlug,
        "authType":     authType,
        "payload":      payload,
    }
    if err := a.client.NerdGraphQueryWithContext(ctx, cloudAuthenticateIntegrationMutation, vars, resp); err != nil {
        return nil, err
    }
    return &resp.CloudAuthenticateIntegration, nil
}
```

### 3.6 New GQL Mutation String — add to `pkg/cloud/cloud_api.go`

```go
const cloudAuthenticateIntegrationMutation = `
mutation(
    $nrAccountId: Int!,
    $providerSlug: CloudProvider!,
    $authType: AuthenticationType!,
    $payload: String!
) {
    cloudAuthenticateIntegration(
        nrAccountId: $nrAccountId
        providerSlug: $providerSlug
        authType: $authType
        payload: $payload
    ) {
        authReferenceId
    }
}`
```

### 3.7 Extend Existing GQL Queries in `pkg/cloud/cloud_api.go`

The `cloudConfigureIntegration` mutation GQL string already includes inline fragments for all existing GCP types. Add inline fragments for `CloudGcpGenericIntegration`:

```graphql
... on CloudGcpGenericIntegration {
    id
    name
    nrAccountId
    metricsPollingInterval
    createdAt
    updatedAt
    service { slug name }
    linkedAccount { id name }
}
```

Add this fragment to:
- `cloudConfigureIntegrationMutation`
- `cloudDisableIntegrationMutation`
- The `GetLinkedAccount` query (inside the `integrations` fragment list)

---

## 4. Resource: `newrelic_cloud_gcp_v2_link_account`

### 4.1 Purpose

Links a GCP project to a New Relic account using Workload Identity Federation. This is a two-step operation executed atomically within the Terraform Create:
1. Call `cloudAuthenticateIntegration` → get `authReferenceId`
2. Call `cloudLinkAccount` with `authReferenceId` → get linked account ID

### 4.2 File

Create: `newrelic/resource_newrelic_cloud_gcp_v2_link_account.go`

### 4.3 Schema

```go
func resourceNewRelicCloudGcpV2LinkAccount() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceNewRelicCloudGcpV2LinkAccountCreate,
        ReadContext:   resourceNewRelicCloudGcpV2LinkAccountRead,
        UpdateContext: resourceNewRelicCloudGcpV2LinkAccountUpdate,
        DeleteContext: resourceNewRelicCloudGcpV2LinkAccountDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: map[string]*schema.Schema{
            "account_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                ForceNew:    true,
                Description: "The New Relic account ID to link the GCP project to.",
            },
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The display name for this linked GCP account in New Relic.",
            },
            "project_id": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The GCP project ID to link (e.g. 'my-gcp-project-123').",
            },
            "wif_credential": {
                Type:      schema.TypeString,
                Required:  true,
                ForceNew:  true,
                Sensitive: true,
                Description: "The Workload Identity Federation credential JSON exported from GCP. " +
                    "Pass the raw JSON string (use file() or jsonencode()). " +
                    "Changing this value forces a new resource (re-link).",
            },
        },
    }
}
```

**Key schema decisions:**
- `wif_credential` is `ForceNew: true` — there is no API to update WIF credentials post-link; a credential change means re-linking.
- `wif_credential` is `Sensitive: true` — prevents credential exposure in plan/apply output and state diff.
- `project_id` is `ForceNew: true` — consistent with v1; changing the project means a new linked account.
- `account_id` is `Computed: true` + `ForceNew: true` — resolved from provider config if not set; once set it cannot change.

### 4.4 CRUD Implementation

#### Create

```go
func resourceNewRelicCloudGcpV2LinkAccountCreate(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    // Step 1: Authenticate (WIF) — get authReferenceId
    wifCredential := d.Get("wif_credential").(string)
    authPayload, err := client.Cloud.CloudAuthenticateIntegrationWithContext(
        ctx,
        accountID,
        "GCP",
        "WIF",
        wifCredential,
    )
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudAuthenticateIntegration failed: %w", err))
    }

    // Step 2: Link account using authReferenceId
    linkInput := cloud.CloudLinkCloudAccountsInput{
        Gcp: []cloud.CloudGcpLinkAccountInput{
            {
                Name:            d.Get("name").(string),
                ProjectId:       d.Get("project_id").(string),
                AuthReferenceId: authPayload.AuthReferenceId,
            },
        },
    }
    linkResp, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountID, linkInput)
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudLinkAccount failed: %w", err))
    }
    if len(linkResp.Errors) > 0 {
        return diag.FromErr(fmt.Errorf("cloudLinkAccount returned errors: %v", linkResp.Errors))
    }
    if len(linkResp.LinkedAccounts) == 0 {
        return diag.FromErr(fmt.Errorf("cloudLinkAccount returned no linked accounts"))
    }

    linkedAccountID := linkResp.LinkedAccounts[0].ID
    d.SetId(strconv.Itoa(linkedAccountID))
    _ = d.Set("account_id", accountID)

    return resourceNewRelicCloudGcpV2LinkAccountRead(ctx, d, meta)
}
```

#### Read

```go
func resourceNewRelicCloudGcpV2LinkAccountRead(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            d.SetId("")
            return nil
        }
        return diag.FromErr(err)
    }

    _ = d.Set("account_id", linkedAccount.NrAccountId)
    _ = d.Set("name", linkedAccount.Name)
    _ = d.Set("project_id", linkedAccount.ExternalId) // GCP project ID stored in ExternalId
    // NOTE: wif_credential is write-only; it is stored in state from the original Create
    // and is NOT populated from the API response (API does not return credentials).

    return nil
}
```

**Important:** `wif_credential` is never set from the Read response because the API does not return credentials. Terraform retains the value written during Create from state. This is the correct and standard pattern for sensitive write-only credentials.

#### Update

Only `name` can be updated in-place. All other fields are `ForceNew`.

```go
func resourceNewRelicCloudGcpV2LinkAccountUpdate(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    renameInput := cloud.CloudRenameAccountsInput{
        LinkedAccountId: linkedAccountID,
        Name:            d.Get("name").(string),
    }
    _, err = client.Cloud.CloudRenameAccountWithContext(ctx, accountID,
        []cloud.CloudRenameAccountsInput{renameInput})
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudRenameAccount failed: %w", err))
    }

    return resourceNewRelicCloudGcpV2LinkAccountRead(ctx, d, meta)
}
```

#### Delete

```go
func resourceNewRelicCloudGcpV2LinkAccountDelete(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    unlinkInput := []cloud.CloudUnlinkAccountsInput{
        {LinkedAccountId: linkedAccountID},
    }
    _, err = client.Cloud.CloudUnlinkAccountWithContext(ctx, accountID, unlinkInput)
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudUnlinkAccount failed: %w", err))
    }

    d.SetId("")
    return nil
}
```

---

## 5. Resource: `newrelic_cloud_gcp_v2_integrations`

### 5.1 Purpose

Configures which GCP services are polled for metrics on a GCP v2 linked account. Supports all 27 existing GCP services plus 7 new GCP v2-only services (34 total). This resource is intended exclusively for accounts created via `newrelic_cloud_gcp_v2_link_account`.

### 5.2 File

Create: `newrelic/resource_newrelic_cloud_gcp_v2_integrations.go`

### 5.3 Complete Service Catalogue

#### Existing 27 services (same API as v1, cleaner HCL names in v2)

| HCL block name | GraphQL field | Extra fields |
|---|---|---|
| `ai_platform` | `gcpAiplatform` | — |
| `alloy_db` | `gcpAlloydb` | — |
| `app_engine` | `gcpAppengine` | — |
| `big_query` | `gcpBigquery` | `fetch_tags`, `fetch_table_metrics` |
| `big_table` | `gcpBigtable` | — |
| `composer` | `gcpComposer` | — |
| `data_flow` | `gcpDataflow` | — |
| `data_proc` | `gcpDataproc` | — |
| `data_store` | `gcpDatastore` | — |
| `firebase_database` | `gcpFirebasedatabase` | — |
| `firebase_hosting` | `gcpFirebasehosting` | — |
| `firebase_storage` | `gcpFirebasestorage` | — |
| `firestore` | `gcpFirestore` | — |
| `functions` | `gcpFunctions` | — |
| `interconnect` | `gcpInterconnect` | — |
| `kubernetes` | `gcpKubernetes` | — |
| `load_balancing` | `gcpLoadbalancing` | — |
| `mem_cache` | `gcpMemcache` | — |
| `pub_sub` | `gcpPubsub` | `fetch_tags` |
| `redis` | `gcpRedis` | — |
| `router` | `gcpRouter` | — |
| `run` | `gcpRun` | — |
| `spanner` | `gcpSpanner` | `fetch_tags` |
| `sql` | `gcpSql` | — |
| `storage` | `gcpStorage` | `fetch_tags` |
| `virtual_machines` | `gcpVms` | — |
| `vpc_access` | `gcpVpcaccess` | — |

#### New 7 GCP v2 services (use `CloudGcpGenericIntegrationInput`)

| HCL block name | GraphQL field | Entity Synthesis |
|---|---|---|
| `api_gateway` | `gcpApiGateway` | Yes |
| `firebase_auth` | `gcpFirebaseAuth` | Yes |
| `firebase_vertex_ai` | `gcpFirebaseVertexAi` | **No** |
| `istio` | `gcpIstio` | **No** |
| `managed_kafka` | `gcpManagedKafka` | Yes |
| `memory_store` | `gcpMemoryStore` | Yes |
| `firebase_app_hosting` | `gcpFirebaseAppHosting` | **No** |

Note: `isEntitySupported = false` for `firebase_vertex_ai`, `istio`, and `firebase_app_hosting`. This is informational; no schema behaviour changes.

### 5.4 Schema

```go
func resourceNewrelicCloudGcpV2Integrations() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceNewrelicCloudGcpV2IntegrationsCreate,
        ReadContext:   resourceNewrelicCloudGcpV2IntegrationsRead,
        UpdateContext: resourceNewrelicCloudGcpV2IntegrationsUpdate,
        DeleteContext: resourceNewrelicCloudGcpV2IntegrationsDelete,
        Importer: &schema.ResourceImporter{
            StateContext: schema.ImportStatePassthroughContext,
        },
        Schema: generateGcpV2IntegrationSchema(),
    }
}
```

#### Schema generation function

```go
func generateGcpV2IntegrationSchema() map[string]*schema.Schema {
    baseSchema := cloudGcpV2IntegrationSchemaBase()
    bigQuerySchema := cloudGcpV2MergeSchema(baseSchema, map[string]*schema.Schema{
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
    fetchTagsSchema := cloudGcpV2MergeSchema(baseSchema, map[string]*schema.Schema{
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
            Description: "The ID of the GCP v2 linked account (from newrelic_cloud_gcp_v2_link_account).",
        },
        // ── Existing 27 services ──
        "ai_platform":        serviceBlock("GCP Vertex AI / AI Platform.", baseSchema),
        "alloy_db":           serviceBlock("GCP AlloyDB.", baseSchema),
        "app_engine":         serviceBlock("GCP App Engine.", baseSchema),
        "big_query":          serviceBlock("GCP BigQuery.", bigQuerySchema),
        "big_table":          serviceBlock("GCP Bigtable.", baseSchema),
        "composer":           serviceBlock("GCP Cloud Composer.", baseSchema),
        "data_flow":          serviceBlock("GCP Cloud Dataflow.", baseSchema),
        "data_proc":          serviceBlock("GCP Cloud Dataproc.", baseSchema),
        "data_store":         serviceBlock("GCP Cloud Datastore.", baseSchema),
        "firebase_database":  serviceBlock("GCP Firebase Realtime Database.", baseSchema),
        "firebase_hosting":   serviceBlock("GCP Firebase Hosting.", baseSchema),
        "firebase_storage":   serviceBlock("GCP Firebase Storage.", baseSchema),
        "firestore":          serviceBlock("GCP Firestore.", baseSchema),
        "functions":          serviceBlock("GCP Cloud Functions.", baseSchema),
        "interconnect":       serviceBlock("GCP Cloud Interconnect.", baseSchema),
        "kubernetes":         serviceBlock("GCP Google Kubernetes Engine (GKE).", baseSchema),
        "load_balancing":     serviceBlock("GCP Cloud Load Balancing.", baseSchema),
        "mem_cache":          serviceBlock("GCP Memcache.", baseSchema),
        "pub_sub":            serviceBlock("GCP Cloud Pub/Sub.", fetchTagsSchema),
        "redis":              serviceBlock("GCP Memorystore for Redis (legacy).", baseSchema),
        "router":             serviceBlock("GCP Cloud Router.", baseSchema),
        "run":                serviceBlock("GCP Cloud Run.", baseSchema),
        "spanner":            serviceBlock("GCP Cloud Spanner.", fetchTagsSchema),
        "sql":                serviceBlock("GCP Cloud SQL.", baseSchema),
        "storage":            serviceBlock("GCP Cloud Storage.", fetchTagsSchema),
        "virtual_machines":   serviceBlock("GCP Compute Engine VMs.", baseSchema),
        "vpc_access":         serviceBlock("GCP Serverless VPC Access.", baseSchema),
        // ── New GCP v2 services ──
        "api_gateway":          serviceBlock("GCP API Gateway (v2 only).", baseSchema),
        "firebase_auth":        serviceBlock("Firebase Authentication (v2 only).", baseSchema),
        "firebase_vertex_ai":   serviceBlock("Firebase Vertex AI (v2 only; no entity synthesis).", baseSchema),
        "istio":                serviceBlock("GCP Istio Service Mesh (v2 only; no entity synthesis).", baseSchema),
        "managed_kafka":        serviceBlock("GCP Managed Service for Apache Kafka (v2 only).", baseSchema),
        "memory_store":         serviceBlock("GCP Memorystore for Redis/Memcached (v2 only).", baseSchema),
        "firebase_app_hosting": serviceBlock("Firebase App Hosting (v2 only; no entity synthesis).", baseSchema),
    }
}

// serviceBlock returns a TypeList schema.Schema with MaxItems:1 for a single integration.
func serviceBlock(description string, elem map[string]*schema.Schema) *schema.Schema {
    return &schema.Schema{
        Type:        schema.TypeList,
        Optional:    true,
        MaxItems:    1,
        Description: description,
        Elem:        &schema.Resource{Schema: elem},
    }
}

// cloudGcpV2IntegrationSchemaBase is the minimal schema all GCP v2 services share.
func cloudGcpV2IntegrationSchemaBase() map[string]*schema.Schema {
    return map[string]*schema.Schema{
        "metrics_polling_interval": {
            Type:        schema.TypeInt,
            Optional:    true,
            Description: "The data polling interval in seconds.",
        },
    }
}

// cloudGcpV2MergeSchema merges two schema maps into a new map (non-destructive).
func cloudGcpV2MergeSchema(base, extra map[string]*schema.Schema) map[string]*schema.Schema {
    result := make(map[string]*schema.Schema, len(base)+len(extra))
    for k, v := range base {
        result[k] = v
    }
    for k, v := range extra {
        result[k] = v
    }
    return result
}
```

### 5.5 CRUD Implementation

#### Create

```go
func resourceNewrelicCloudGcpV2IntegrationsCreate(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID := d.Get("linked_account_id").(int)
    configureInput, _ := expandCloudGcpV2IntegrationsInput(d, linkedAccountID)

    resp, err := client.Cloud.CloudConfigureIntegrationWithContext(
        ctx, accountID, configureInput,
    )
    if err != nil {
        return diag.FromErr(err)
    }
    if len(resp.Errors) > 0 {
        return diag.FromErr(fmt.Errorf("cloudConfigureIntegration errors: %v", resp.Errors))
    }

    d.SetId(strconv.Itoa(linkedAccountID))
    _ = d.Set("account_id", accountID)

    return resourceNewrelicCloudGcpV2IntegrationsRead(ctx, d, meta)
}
```

#### Read

```go
func resourceNewrelicCloudGcpV2IntegrationsRead(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            d.SetId("")
            return nil
        }
        return diag.FromErr(err)
    }

    _ = d.Set("account_id", linkedAccount.NrAccountId)
    _ = d.Set("linked_account_id", linkedAccountID)

    return flattenCloudGcpV2LinkedAccount(d, linkedAccount)
}
```

#### Update

```go
func resourceNewrelicCloudGcpV2IntegrationsUpdate(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    configureInput, disableInput := expandCloudGcpV2IntegrationsInput(d, linkedAccountID)

    // Disable removed integrations first
    _, err = client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudDisableIntegration failed: %w", err))
    }

    // Enable/update requested integrations
    resp, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)
    if err != nil {
        return diag.FromErr(fmt.Errorf("cloudConfigureIntegration failed: %w", err))
    }
    if len(resp.Errors) > 0 {
        return diag.FromErr(fmt.Errorf("cloudConfigureIntegration errors: %v", resp.Errors))
    }

    return resourceNewrelicCloudGcpV2IntegrationsRead(ctx, d, meta)
}
```

#### Delete

```go
func resourceNewrelicCloudGcpV2IntegrationsDelete(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)

    linkedAccountID, err := strconv.Atoi(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    _, disableInput := expandCloudGcpV2IntegrationsInput(d, linkedAccountID)
    _, err = client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId("")
    return nil
}
```

### 5.6 Expand / Flatten Helper Functions

#### expandCloudGcpV2IntegrationsInput

Returns both the configure input and disable input from current schema state.

```go
// expandCloudGcpV2IntegrationsInput builds CloudIntegrationsInput (enable) and
// CloudDisableIntegrationsInput (disable) from the resource schema.
// Both are always built — callers use configure for present blocks and disable for absent ones.
func expandCloudGcpV2IntegrationsInput(
    d *schema.ResourceData,
    linkedAccountID int,
) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {

    gcpInput := cloud.CloudGcpIntegrationsInput{}
    gcpDisable := cloud.CloudGcpDisableIntegrationsInput{}
    disable := cloud.CloudDisableAccountIntegrationInput{LinkedAccountId: linkedAccountID}

    // ── Helper closures ──────────────────────────────────────────────────────
    present := func(key string) bool {
        v, ok := d.GetOk(key)
        if !ok {
            return false
        }
        l, ok := v.([]interface{})
        return ok && len(l) > 0
    }

    baseInput := func(key string) cloud.CloudGcpAlloydbIntegrationInput {
        // generic helper for services with only metricsPollingInterval
        in := cloud.CloudGcpAlloydbIntegrationInput{LinkedAccountId: linkedAccountID}
        if v := d.Get(key + ".0.metrics_polling_interval"); v != nil {
            in.MetricsPollingInterval = v.(int)
        }
        return in
    }
    // (inline definition above is illustrative — use the actual service type per field below)

    genericInput := func(key string) cloud.CloudGcpGenericIntegrationInput {
        in := cloud.CloudGcpGenericIntegrationInput{LinkedAccountId: linkedAccountID}
        if v, ok := d.GetOk(key + ".0.metrics_polling_interval"); ok {
            in.MetricsPollingInterval = v.(int)
        }
        return in
    }

    // ── Existing 27 services ─────────────────────────────────────────────────
    // Each block: if present → add to configure input; else → add to disable input

    if present("ai_platform") {
        gcpInput.GcpAiplatform = []cloud.CloudGcpAiplatformIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("ai_platform.0.metrics_polling_interval").(int),
        }}
    } else {
        gcpDisable.GcpAiplatform = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    if present("alloy_db") {
        gcpInput.GcpAlloydb = []cloud.CloudGcpAlloydbIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("alloy_db.0.metrics_polling_interval").(int),
        }}
    } else {
        gcpDisable.GcpAlloydb = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    if present("app_engine") {
        gcpInput.GcpAppengine = []cloud.CloudGcpAppengineIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("app_engine.0.metrics_polling_interval").(int),
        }}
    } else {
        gcpDisable.GcpAppengine = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    if present("big_query") {
        gcpInput.GcpBigquery = []cloud.CloudGcpBigqueryIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("big_query.0.metrics_polling_interval").(int),
            FetchTags:              d.Get("big_query.0.fetch_tags").(bool),
            FetchTableMetrics:      d.Get("big_query.0.fetch_table_metrics").(bool),
        }}
    } else {
        gcpDisable.GcpBigquery = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    if present("big_table") {
        gcpInput.GcpBigtable = []cloud.CloudGcpBigtableIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("big_table.0.metrics_polling_interval").(int),
        }}
    } else {
        gcpDisable.GcpBigtable = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    // ... repeat the same pattern for the remaining 22 existing services:
    // composer, data_flow, data_proc, data_store, firebase_database, firebase_hosting,
    // firebase_storage, firestore, functions, interconnect, kubernetes, load_balancing,
    // mem_cache, pub_sub (with fetch_tags), redis, router, run,
    // spanner (with fetch_tags), sql, storage (with fetch_tags), virtual_machines, vpc_access

    // pub_sub example (with fetch_tags):
    if present("pub_sub") {
        gcpInput.GcpPubsub = []cloud.CloudGcpPubsubIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("pub_sub.0.metrics_polling_interval").(int),
            FetchTags:              d.Get("pub_sub.0.fetch_tags").(bool),
        }}
    } else {
        gcpDisable.GcpPubsub = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    // spanner example (with fetch_tags):
    if present("spanner") {
        gcpInput.GcpSpanner = []cloud.CloudGcpSpannerIntegrationInput{{
            LinkedAccountId:        linkedAccountID,
            MetricsPollingInterval: d.Get("spanner.0.metrics_polling_interval").(int),
            FetchTags:              d.Get("spanner.0.fetch_tags").(bool),
        }}
    } else {
        gcpDisable.GcpSpanner = []cloud.CloudDisableAccountIntegrationInput{disable}
    }

    // ── New GCP v2 services (all use CloudGcpGenericIntegrationInput) ─────────

    v2Services := []struct {
        key      string
        setFn    func([]cloud.CloudGcpGenericIntegrationInput)
        disableFn func([]cloud.CloudDisableAccountIntegrationInput)
    }{
        {"api_gateway",          func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpApiGateway = v },          func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpApiGateway = v }},
        {"firebase_auth",        func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpFirebaseAuth = v },        func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpFirebaseAuth = v }},
        {"firebase_vertex_ai",   func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpFirebaseVertexAi = v },   func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpFirebaseVertexAi = v }},
        {"istio",                func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpIstio = v },              func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpIstio = v }},
        {"managed_kafka",        func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpManagedKafka = v },       func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpManagedKafka = v }},
        {"memory_store",         func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpMemoryStore = v },        func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpMemoryStore = v }},
        {"firebase_app_hosting", func(v []cloud.CloudGcpGenericIntegrationInput) { gcpInput.GcpFirebaseAppHosting = v }, func(v []cloud.CloudDisableAccountIntegrationInput) { gcpDisable.GcpFirebaseAppHosting = v }},
    }

    for _, svc := range v2Services {
        if present(svc.key) {
            svc.setFn([]cloud.CloudGcpGenericIntegrationInput{genericInput(svc.key)})
        } else {
            svc.disableFn([]cloud.CloudDisableAccountIntegrationInput{disable})
        }
    }

    return cloud.CloudIntegrationsInput{Gcp: gcpInput},
           cloud.CloudDisableIntegrationsInput{Gcp: gcpDisable}
}
```

#### flattenCloudGcpV2LinkedAccount

```go
func flattenCloudGcpV2LinkedAccount(
    d *schema.ResourceData,
    linkedAccount *cloud.CloudLinkedAccount,
) diag.Diagnostics {
    for _, rawIntegration := range linkedAccount.Integrations {
        switch v := rawIntegration.(type) {
        case *cloud.CloudGcpAiplatformIntegration:
            _ = d.Set("ai_platform", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpAlloydbIntegration:
            _ = d.Set("alloy_db", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpAppengineIntegration:
            _ = d.Set("app_engine", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpBigqueryIntegration:
            _ = d.Set("big_query", []interface{}{map[string]interface{}{
                "metrics_polling_interval": v.MetricsPollingInterval,
                "fetch_tags":              v.FetchTags,
                "fetch_table_metrics":     v.FetchTableMetrics,
            }})
        case *cloud.CloudGcpBigtableIntegration:
            _ = d.Set("big_table", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpComposerIntegration:
            _ = d.Set("composer", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpDataflowIntegration:
            _ = d.Set("data_flow", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpDataprocIntegration:
            _ = d.Set("data_proc", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpDatastoreIntegration:
            _ = d.Set("data_store", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpFirebasedatabaseIntegration:
            _ = d.Set("firebase_database", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpFirebasehostingIntegration:
            _ = d.Set("firebase_hosting", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpFirebasestorageIntegration:
            _ = d.Set("firebase_storage", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpFirestoreIntegration:
            _ = d.Set("firestore", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpFunctionsIntegration:
            _ = d.Set("functions", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpInterconnectIntegration:
            _ = d.Set("interconnect", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpKubernetesIntegration:
            _ = d.Set("kubernetes", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpLoadbalancingIntegration:
            _ = d.Set("load_balancing", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpMemcacheIntegration:
            _ = d.Set("mem_cache", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpPubsubIntegration:
            _ = d.Set("pub_sub", []interface{}{map[string]interface{}{
                "metrics_polling_interval": v.MetricsPollingInterval,
                "fetch_tags":              v.FetchTags,
            }})
        case *cloud.CloudGcpRedisIntegration:
            _ = d.Set("redis", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpRouterIntegration:
            _ = d.Set("router", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpRunIntegration:
            _ = d.Set("run", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpSpannerIntegration:
            _ = d.Set("spanner", []interface{}{map[string]interface{}{
                "metrics_polling_interval": v.MetricsPollingInterval,
                "fetch_tags":              v.FetchTags,
            }})
        case *cloud.CloudGcpSqlIntegration:
            _ = d.Set("sql", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpStorageIntegration:
            _ = d.Set("storage", []interface{}{map[string]interface{}{
                "metrics_polling_interval": v.MetricsPollingInterval,
                "fetch_tags":              v.FetchTags,
            }})
        case *cloud.CloudGcpVmsIntegration:
            _ = d.Set("virtual_machines", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))
        case *cloud.CloudGcpVpcaccessIntegration:
            _ = d.Set("vpc_access", flattenGcpV2BaseIntegration(v.MetricsPollingInterval))

        // New GCP v2 services — all return CloudGcpGenericIntegration.
        // Distinguished by service.slug since they share the same Go type.
        case *cloud.CloudGcpGenericIntegration:
            block := flattenGcpV2BaseIntegration(v.MetricsPollingInterval)
            switch v.Service.Slug {
            case "gcpApiGateway":
                _ = d.Set("api_gateway", block)
            case "gcpFirebaseAuth":
                _ = d.Set("firebase_auth", block)
            case "gcpFirebaseVertexAi":
                _ = d.Set("firebase_vertex_ai", block)
            case "gcpIstio":
                _ = d.Set("istio", block)
            case "gcpManagedKafka":
                _ = d.Set("managed_kafka", block)
            case "gcpMemoryStore":
                _ = d.Set("memory_store", block)
            case "gcpFirebaseAppHosting":
                _ = d.Set("firebase_app_hosting", block)
            }
        }
    }
    return nil
}

func flattenGcpV2BaseIntegration(pollingInterval int) []interface{} {
    return []interface{}{map[string]interface{}{
        "metrics_polling_interval": pollingInterval,
    }}
}
```

---

## 6. Data Source: `newrelic_cloud_gcp_v2_account`

### 6.1 Purpose

Enables users to look up a GCP v2 linked account by name to use its ID in other resources (e.g., referencing an account created outside Terraform).

### 6.2 File

Create: `newrelic/data_source_newrelic_cloud_gcp_v2_account.go`

### 6.3 Schema and Implementation

```go
func dataSourceNewRelicCloudGcpV2Account() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceNewRelicCloudGcpV2AccountRead,
        Schema: map[string]*schema.Schema{
            "account_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "The New Relic account ID. Defaults to the provider account.",
            },
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The name of the GCP v2 linked account to look up.",
            },
        },
    }
}

func dataSourceNewRelicCloudGcpV2AccountRead(
    ctx context.Context,
    d *schema.ResourceData,
    meta interface{},
) diag.Diagnostics {
    providerConfig := meta.(*ProviderConfig)
    client := providerConfig.NewClient
    accountID := selectAccountID(providerConfig, d)
    targetName := d.Get("name").(string)

    // GCP v2 linked accounts use the "gcp_v2" provider slug internally.
    // The public API exposes them under "gcp_v2" in GetLinkedAccounts.
    // If "gcp_v2" returns empty, fall back to "gcp" (verify with API team which slug works).
    linkedAccounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, "gcp_v2")
    if err != nil {
        return diag.FromErr(fmt.Errorf("GetLinkedAccounts failed: %w", err))
    }

    for _, account := range *linkedAccounts {
        if strings.EqualFold(account.Name, targetName) && account.NrAccountId == accountID {
            d.SetId(strconv.Itoa(account.ID))
            _ = d.Set("account_id", account.NrAccountId)
            _ = d.Set("name", account.Name)
            return nil
        }
    }

    return diag.Errorf(
        "no GCP v2 linked account named %q found for New Relic account %d",
        targetName, accountID,
    )
}
```

**Implementation note:** The `GetLinkedAccountsWithContext` call uses `"gcp_v2"` as the provider string. This needs to be validated against the live API during testing. If the API does not expose `gcp_v2` to callers, use `"gcp"` and add a comment explaining the fallback.

---

## 7. File Structure

### Files to Create

```
newrelic/
├── resource_newrelic_cloud_gcp_v2_link_account.go
├── resource_newrelic_cloud_gcp_v2_link_account_test.go
├── resource_newrelic_cloud_gcp_v2_integrations.go
├── resource_newrelic_cloud_gcp_v2_integrations_test.go
└── data_source_newrelic_cloud_gcp_v2_account.go
```

### Files to Modify

```
newrelic/provider_newrelic.go         — register new resources and data source
go.mod / go.sum                       — updated newrelic-client-go version
```

### Files NOT Modified (guaranteed backward compatibility)

```
newrelic/resource_newrelic_cloud_gcp_integrations.go      — unchanged
newrelic/resource_newrelic_cloud_gcp_link_account.go      — unchanged
newrelic/data_source_newrelic_cloud_account.go            — unchanged
```

---

## 8. Provider Registration

Add to the `ResourcesMap` and `DataSourcesMap` in `newrelic/provider_newrelic.go`:

```go
// In ResourcesMap:
"newrelic_cloud_gcp_v2_link_account":   resourceNewRelicCloudGcpV2LinkAccount(),
"newrelic_cloud_gcp_v2_integrations":   resourceNewrelicCloudGcpV2Integrations(),

// In DataSourcesMap:
"newrelic_cloud_gcp_v2_account": dataSourceNewRelicCloudGcpV2Account(),
```

---

## 9. Error Handling Requirements

1. **cloudAuthenticateIntegration failure**: Wrap error with context — `"cloudAuthenticateIntegration failed: %w"`. Do not retry; the WIF JSON is invalid or the session store is unavailable.

2. **authReferenceId TTL**: The session is valid for 30 minutes. The two-step Create (authenticate + link) runs sequentially within a single Terraform apply operation — TTL expiry is not a concern in practice. No retry logic needed.

3. **cloudLinkAccount errors array**: Always check `linkResp.Errors`. Return the errors as a diagnostic even if `err == nil`.

4. **cloudConfigureIntegration errors array**: Same pattern — check `resp.Errors` in addition to `err`.

5. **Not found on Read**: Detect by checking `err.Error()` for `"not found"` (same pattern as v1). Clear `d.SetId("")` and return `nil` to let Terraform detect drift.

6. **Import**: All three resources support `schema.ImportStatePassthroughContext`. The resource ID is the string representation of the linked account ID integer.

---

## 10. Testing Requirements

### 10.1 Build tags

All test files must include:
```go
//go:build integration || CLOUD
```

### 10.2 Required environment variables

| Variable | Used by |
|---|---|
| `NEW_RELIC_SUBACCOUNT_ID` | Both resource tests |
| `INTEGRATION_TESTING_GCP_PROJECT_ID` | Both resource tests |
| `INTEGRATION_TESTING_GCP_WIF_CREDENTIAL` | Link account test |

`INTEGRATION_TESTING_GCP_WIF_CREDENTIAL` should contain the full WIF JSON string (e.g., `$(cat wif-credential.json)`).

### 10.3 Link account test (`resource_newrelic_cloud_gcp_v2_link_account_test.go`)

```go
func TestAccNewRelicCloudGcpV2LinkAccount_Basic(t *testing.T) {
    testProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
    testWifCredential := os.Getenv("INTEGRATION_TESTING_GCP_WIF_CREDENTIAL")
    if testProjectID == "" || testWifCredential == "" {
        t.Skip("skipping: INTEGRATION_TESTING_GCP_PROJECT_ID and INTEGRATION_TESTING_GCP_WIF_CREDENTIAL required")
    }

    resource.ParallelTest(t, resource.TestCase{
        PreCheck:          func() { testAccPreCheck(t) },
        ProviderFactories: providerFactories,
        CheckDestroy:      testAccCheckNewRelicCloudGcpV2LinkAccountDestroyed,
        Steps: []resource.TestStep{
            // Create
            {
                Config: testAccNewRelicCloudGcpV2LinkAccountConfig(testProjectID, testWifCredential, "tf-test-gcp-v2"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("newrelic_cloud_gcp_v2_link_account.test", "name", "tf-test-gcp-v2"),
                    resource.TestCheckResourceAttr("newrelic_cloud_gcp_v2_link_account.test", "project_id", testProjectID),
                    resource.TestCheckResourceAttrSet("newrelic_cloud_gcp_v2_link_account.test", "id"),
                ),
            },
            // Rename (update name only)
            {
                Config: testAccNewRelicCloudGcpV2LinkAccountConfig(testProjectID, testWifCredential, "tf-test-gcp-v2-renamed"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("newrelic_cloud_gcp_v2_link_account.test", "name", "tf-test-gcp-v2-renamed"),
                ),
            },
            // Import
            {
                ResourceName:            "newrelic_cloud_gcp_v2_link_account.test",
                ImportState:             true,
                ImportStateVerify:       true,
                ImportStateVerifyIgnore: []string{"wif_credential"}, // write-only, not in API response
            },
        },
    })
}
```

Note: `wif_credential` must be in `ImportStateVerifyIgnore` because the API does not return credentials.

### 10.4 Integrations test (`resource_newrelic_cloud_gcp_v2_integrations_test.go`)

```go
func TestAccNewRelicCloudGcpV2Integrations_Basic(t *testing.T) {
    // Test steps:
    // Step 1: Create link account + configure big_query + api_gateway
    // Step 2: Update — add firebase_auth, remove api_gateway
    // Step 3: Import state
    // Destroy: delete integrations, then link account
}
```

The test config should exercise:
- At least one existing service (e.g., `big_query` with `fetch_tags` and `fetch_table_metrics`)
- At least one new GCP v2 service (e.g., `api_gateway`)
- An update that adds and removes services

### 10.5 Data source test

```go
func TestAccNewRelicCloudGcpV2Account_DataSource(t *testing.T) {
    // Create a link account resource, then verify the data source finds it by name
}
```

---

## 11. Complete Example Usage (for documentation)

### Minimal — link account and enable two services

```hcl
resource "newrelic_cloud_gcp_v2_link_account" "example" {
  account_id     = 12345678
  name           = "production-gcp"
  project_id     = "my-gcp-project-id"
  wif_credential = file("${path.module}/wif-credential.json")
}

resource "newrelic_cloud_gcp_v2_integrations" "example" {
  account_id        = newrelic_cloud_gcp_v2_link_account.example.account_id
  linked_account_id = newrelic_cloud_gcp_v2_link_account.example.id

  big_query {
    metrics_polling_interval = 300
    fetch_tags               = true
    fetch_table_metrics      = true
  }

  api_gateway {
    metrics_polling_interval = 300
  }
}
```

### Full — all 7 new GCP v2 services

```hcl
resource "newrelic_cloud_gcp_v2_integrations" "full" {
  account_id        = newrelic_cloud_gcp_v2_link_account.example.account_id
  linked_account_id = newrelic_cloud_gcp_v2_link_account.example.id

  # GCP v2-only services
  api_gateway {
    metrics_polling_interval = 300
  }
  firebase_auth {
    metrics_polling_interval = 300
  }
  firebase_vertex_ai {
    metrics_polling_interval = 300
  }
  istio {
    metrics_polling_interval = 300
  }
  managed_kafka {
    metrics_polling_interval = 300
  }
  memory_store {
    metrics_polling_interval = 300
  }
  firebase_app_hosting {
    metrics_polling_interval = 300
  }

  # Also works with existing GCP services
  cloud_run:
  run {
    metrics_polling_interval = 60
  }
  kubernetes {
    metrics_polling_interval = 60
  }
}
```

### Data source usage

```hcl
# Look up an existing GCP v2 account (created outside Terraform or in another state)
data "newrelic_cloud_gcp_v2_account" "existing" {
  account_id = 12345678
  name       = "production-gcp"
}

resource "newrelic_cloud_gcp_v2_integrations" "external" {
  account_id        = data.newrelic_cloud_gcp_v2_account.existing.account_id
  linked_account_id = data.newrelic_cloud_gcp_v2_account.existing.id

  big_query {
    metrics_polling_interval = 300
  }
}
```

---

## 12. Backward Compatibility Guarantees

| What | Why it is safe |
|---|---|
| `newrelic_cloud_gcp_link_account` | Not touched. Zero changes. |
| `newrelic_cloud_gcp_integrations` | Not touched. Zero changes. |
| `newrelic_cloud_account` data source | Not touched. Zero changes. |
| `CloudGcpLinkAccountInput.AuthReferenceId` | `omitempty` ensures existing v1 serialisations remain identical. |
| 7 new fields in `CloudGcpIntegrationsInput` | All `omitempty`; existing configure calls with 0 of these fields serialise identically. |
| 7 new fields in `CloudGcpDisableIntegrationsInput` | Same as above. |
| Provider registration | Additions only; no existing key renamed or removed. |

---

## 13. Open Questions for Implementing Engineer

1. **`GetLinkedAccountsWithContext("gcp_v2")` vs `"gcp"`**: Verify at the API level whether GCP v2 linked accounts are returned when querying with provider string `"gcp_v2"` or `"gcp"`. Update `data_source_newrelic_cloud_gcp_v2_account.go` accordingly.

2. **Service slug casing for `CloudGcpGenericIntegration`**: Confirm the exact `Service.Slug` values returned by the API for the 7 new services (e.g., `"gcpApiGateway"` vs `"gcp_api_gateway"`). The flatten switch in Section 5.6 uses `"gcpApiGateway"` per the GraphQL field name convention — validate this against a live API response.

3. **`cloudAuthenticateIntegration` enum serialisation**: Confirm whether the `CloudProvider` and `AuthenticationType` enums are sent as uppercase strings (`"GCP"`, `"WIF"`) or lowercase. The GraphQL schema shows `GCP` and `WIF` as enum values; use those exact strings.

4. **`ai_platform` in v1**: The existing `newrelic_cloud_gcp_integrations` (v1) does **not** have an `ai_platform` block despite `GcpAiplatformIntegrationInput` existing in client-go. The v2 integrations resource adds it. Confirm the API accepts `gcpAiplatform` for v2 linked accounts.

---

## 14. Summary Checklist for Implementing Engineer

- [ ] Apply Section 3 changes to `newrelic-client-go` (types.go, cloud_api.go, cloud.go)
- [ ] Update `go.mod` in Terraform provider to reference new client-go version
- [ ] Create `resource_newrelic_cloud_gcp_v2_link_account.go` per Section 4
- [ ] Create `resource_newrelic_cloud_gcp_v2_integrations.go` per Section 5 (implement all 27 existing services in `expandCloudGcpV2IntegrationsInput` — Section 5.6 gives the complete pattern, implement remaining 22 services following the same template)
- [ ] Create `data_source_newrelic_cloud_gcp_v2_account.go` per Section 6
- [ ] Register all three in `provider_newrelic.go` per Section 8
- [ ] Resolve open questions (Section 13) before writing tests
- [ ] Create test files per Section 10
- [ ] Run `go vet ./...` and `go build ./...` — zero errors required
- [ ] Verify `newrelic_cloud_gcp_link_account` and `newrelic_cloud_gcp_integrations` acceptance tests still pass unchanged
