# New Relic Notebooks – Blob Storage API Reference

A practical guide for manually testing every CRUD operation against the
Notebooks Blob Storage API. All examples use `curl`. All operations require a
New Relic **User API key** (starts `NRAK-`).

---

## Prerequisites

Set these shell variables once; every example below references them:

```bash
export NR_API_KEY="NRAK-<your-user-api-key>"
export NR_ORG_ID="<your-org-uuid>"          # see "Get your Org ID" below
export NR_ACCOUNT_ID=<your-account-id>       # numeric, for NRQL widgets
export BLOB_BASE="https://blob-api.service.newrelic.com/v1/e"
```

For EU accounts, use `https://blob-api.service.eu.newrelic.com/v1/e`.
For JP accounts, use `https://blob-api.service.jp.newrelic.com/v1/e`.

### Get your Org ID

Your Org ID is a UUID. Retrieve it from NerdGraph:

```bash
curl -s -X POST https://api.newrelic.com/graphql \
  -H "API-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  --data '{"query":"{ actor { organization { id } } }"}' \
  | jq -r '.data.actor.organization.id'
```

---

## Headers reference

| Header | Required | Value |
|---|---|---|
| `Api-Key` | All requests | Your NRAK-... User API key |
| `Content-Type` | POST only | `application/json` |
| `NewRelic-Entity` | Create & Rename | JSON string: `{"name":"<title>"}` |

---

## Create a notebook

```bash
curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "my-test-notebook"}' \
  -d '{
    "version": "1",
    "blocks": [
      {
        "type": "widget",
        "content": {
          "type": "visualization",
          "id": "viz.markdown",
          "props": {
            "text": "# Investigation notes\n\nAdd context here."
          }
        }
      },
      {
        "type": "widget",
        "content": {
          "type": "visualization",
          "id": "viz.billboard",
          "props": {
            "nrqlQueries": [
              {
                "accountIds": ['"$NR_ACCOUNT_ID"'],
                "query": "SELECT count(*) FROM Transaction SINCE 1 hour ago"
              }
            ]
          }
        }
      }
    ]
  }' | jq .
```

**Save the `entityGuid`** from the response — it is the stable notebook ID for
all subsequent calls:

```bash
GUID=$(curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "my-test-notebook"}' \
  -d '{"version":"1","blocks":[]}' | jq -r '.entityGuid')
echo "GUID: $GUID"
```

**Expected response (HTTP 200):**
```json
{
  "entityGuid": "NjQyNTg2NX...",
  "blobId":     "MXxFTlRJVF...",
  "blobVersionEntity": null
}
```

---

## Read notebook content

```bash
curl -s -X GET "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" | jq .
```

Returns the raw JSON body that was last written — no envelope, no metadata.
Returns `"Blob not found."` with HTTP 404 if the notebook was deleted.

---

## Read notebook metadata (NerdGraph)

Use NerdGraph to retrieve the title, tags, scope, and version counter:

```bash
curl -s -X POST https://api.newrelic.com/graphql \
  -H "API-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  --data "{
    \"query\": \"query(\\\$id: ID!) { actor { entityManagement { entity(id: \\\$id) { __typename ... on EntityManagementNotebookEntity { id name type scope { id type } metadata { version createdAt updatedAt } } } } } }\",
    \"variables\": {\"id\": \"$GUID\"}
  }" | jq .
```

---

## Update notebook content

POST to the same URL as Read (using the entity GUID):

```bash
curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1",
    "blocks": [
      {
        "type": "widget",
        "content": {
          "type": "visualization",
          "id": "viz.markdown",
          "props": {
            "text": "# Updated notes\n\nContent changed successfully."
          }
        }
      },
      {
        "type": "widget",
        "content": {
          "type": "visualization",
          "id": "viz.line",
          "props": {
            "nrqlQueries": [
              {
                "accountIds": ['"$NR_ACCOUNT_ID"'],
                "query": "SELECT count(*) FROM Transaction TIMESERIES"
              }
            ]
          }
        }
      }
    ]
  }' | jq .
```

Each POST to the update URL creates a new immutable revision. The
`metadata.version` on the NerdGraph entity increments with each write
(1 on create, 2 after the first update, and so on).

---

## Rename a notebook (atomic with content update)

The Blob Storage API has **no rename-only path**. Supply a `NewRelic-Entity`
header with the new name on any content POST and the rename is atomic:

```bash
curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "my-test-notebook-renamed"}' \
  -d '{"version":"1","blocks":[]}' | jq .
```

You must provide the current (or new) content body even if you only want
to rename. Fetch the current content first if you want a rename-only effect:

```bash
# Rename-only: fetch current content then re-POST with new name
CURRENT=$(curl -s "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY")

curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "my-test-notebook-renamed"}' \
  -d "$CURRENT" | jq .
```

---

## Delete a notebook

```bash
curl -s -X DELETE "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -w "\nHTTP %{http_code}\n"
```

Returns **HTTP 200** with an empty body on success (the documentation says 204
but the live API returns 200). After deletion, a GET returns HTTP 404
`"Blob not found."`.

---

## List all notebooks (NerdGraph)

NerdGraph's `entitySearch` with the `NOTEBOOK` type filter lists all notebooks
visible to the authenticated user in their default scope:

```bash
curl -s -X POST https://api.newrelic.com/graphql \
  -H "API-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  --data '{
    "query": "query { actor { entityManagement { entitySearch(query: \"type = '"'"'NOTEBOOK'"'"'\") { entities { __typename ... on EntityManagementNotebookEntity { id name scope { id type } metadata { version updatedAt } } } nextCursor } } } }"
  }' | jq .
```

Narrow to a specific org:

```bash
QUERY="type = 'NOTEBOOK' AND scope.id = '$NR_ORG_ID'"
# URL-encode and pass as a variable
curl -s -X POST https://api.newrelic.com/graphql \
  -H "API-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  --data "{\"query\":\"query(\\\$q: String!) { actor { entityManagement { entitySearch(query: \\\$q) { entities { __typename ... on EntityManagementNotebookEntity { id name } } } } } }\", \"variables\":{\"q\":\"$QUERY\"}}" \
  | jq .
```

---

## Notebook content schema

The Blob API stores arbitrary versioned JSON. The minimum viable body is:

```json
{ "version": "1", "blocks": [] }
```

Each element in `blocks` is a widget or container:

### Markdown widget
```json
{
  "type": "widget",
  "content": {
    "type": "visualization",
    "id": "viz.markdown",
    "props": {
      "text": "# Heading\n\nMarkdown content here."
    }
  }
}
```

### Billboard widget
```json
{
  "type": "widget",
  "content": {
    "type": "visualization",
    "id": "viz.billboard",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction SINCE 1 hour ago"
        }
      ],
      "thresholdsWithSeriesOverrides": {
        "thresholds": [
          { "to": 1000, "severity": "success" },
          { "from": 1000, "to": 5000, "severity": "warning" },
          { "from": 5000, "severity": "critical" }
        ]
      }
    }
  }
}
```

### Line chart widget
```json
{
  "type": "widget",
  "content": {
    "type": "visualization",
    "id": "viz.line",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction TIMESERIES AUTO"
        }
      ]
    }
  }
}
```

### Container (groups widgets with a layout)
```json
{
  "type": "container",
  "props": { "layout": "stack" },
  "content": [
    { "type": "widget", "content": { "..." } },
    { "type": "widget", "content": { "..." } }
  ]
}
```

---

## Common widget `id` values

| `id` | Widget type |
|---|---|
| `viz.markdown` | Static text / markdown |
| `viz.billboard` | Single metric with thresholds |
| `viz.line` | Time-series line chart |
| `viz.area` | Time-series area chart |
| `viz.bar` | Bar chart |
| `viz.pie` | Pie chart |
| `viz.table` | Data table |
| `viz.histogram` | Histogram |
| `viz.heatmap` | Heatmap |
| `viz.stacked-bar` | Stacked bar chart |
| `viz.json` | Raw JSON display |

---

## Error reference

| HTTP status | Body | Meaning |
|---|---|---|
| 200 | `{entityGuid, blobId, ...}` | Create / update succeeded |
| 200 | _(empty)_ | Delete succeeded |
| 400 | _(text)_ | Bad request — often a duplicate name (names must be unique per org) or malformed JSON |
| 401 | `Unauthorized` | Invalid or missing API key |
| 404 | `"Blob not found."` | Notebook does not exist or was deleted |

---

## Quick end-to-end smoke test

Run this entire block to create, verify, rename, and delete a notebook in one pass:

```bash
export NR_API_KEY="NRAK-<your-key>"
export NR_ORG_ID="<your-org-uuid>"
export BLOB_BASE="https://blob-api.service.newrelic.com/v1/e"

echo "=== CREATE ==="
GUID=$(curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "smoke-test-notebook"}' \
  -d '{"version":"1","blocks":[{"type":"widget","content":{"type":"visualization","id":"viz.markdown","props":{"text":"smoke test"}}}]}' \
  | jq -r '.entityGuid')
echo "GUID: $GUID"

echo ""
echo "=== READ ==="
curl -s "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" | jq .

echo ""
echo "=== RENAME ==="
curl -s -X POST "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -H "Content-Type: application/json" \
  -H 'NewRelic-Entity: {"name": "smoke-test-notebook-renamed"}' \
  -d '{"version":"1","blocks":[{"type":"widget","content":{"type":"visualization","id":"viz.markdown","props":{"text":"renamed"}}}]}' \
  | jq -r '.entityGuid'

echo ""
echo "=== DELETE ==="
curl -s -X DELETE "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -w "HTTP %{http_code}\n"

echo ""
echo "=== VERIFY GONE (expect 404) ==="
curl -s "$BLOB_BASE/organizations/$NR_ORG_ID/Notebooks/$GUID" \
  -H "Api-Key: $NR_API_KEY" \
  -w " (HTTP %{http_code})\n"
```
