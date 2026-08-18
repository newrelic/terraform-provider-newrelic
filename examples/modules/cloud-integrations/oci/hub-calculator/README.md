# OCI Connector Hub Calculator

Computes the `connector_hubs_data` variable for the [`metrics-integration`](../metrics-integration) Terraform module.

OCI caps each Service Connector Hub monitoring source at **5 compartments** and **50 namespaces total**. For tenancies with many compartments or dense namespace selections, calculating this by hand is error-prone. This script replicates the same bin-packing algorithm used by the New Relic backend and outputs a value you can paste directly into your `.tfvars`.

## Requirements

Python 3.9 or later. No third-party dependencies.

## Usage

```bash
python3 connector_hub_calculator.py --input compartments.json
```

Pipe the output directly into your `.tfvars`:

```bash
python3 connector_hub_calculator.py --input compartments.json >> terraform.tfvars
```

Use `--pretty` to inspect the hub layout before using it:

```bash
python3 connector_hub_calculator.py --input compartments.json --pretty
```

## Input format

Create a JSON file describing your region and the compartments you want to instrument:

```json
{
  "region": "us-ashburn-1",
  "compartments": [
    {
      "compartment_id": "ocid1.compartment.oc1..aaaa",
      "namespaces": ["oci_compute", "oci_lbaas", "oci_faas"]
    },
    {
      "compartment_id": "ocid1.compartment.oc1..bbbb",
      "namespaces": ["oci_compute", "oci_blockstore"]
    }
  ]
}
```

| Field | Type | Description |
|---|---|---|
| `region` | string | OCI region identifier (e.g. `us-ashburn-1`) — used to generate hub names |
| `compartments` | array | Compartments to instrument |
| `compartments[].compartment_id` | string | Compartment OCID |
| `compartments[].namespaces` | array of strings | OCI metric namespaces to collect for this compartment |

The `name` and `description` fields in each hub are generated automatically — you do not supply them.

See [Supported OCI service categories](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#supported-oci-service-categories) for the full namespace list.

## Output

Without `--pretty`, the script prints a single-line quoted JSON string ready to assign to `connector_hubs_data`:

```
"[{\"name\":\"newrelic-metrics-connector-hub-us-ashburn-1\",\"description\":\"[DO NOT DELETE] New Relic Metrics Connector Hub\",\"compartments\":[...]}]"
```

Copy this value into your `.tfvars`:

```hcl
connector_hubs_data = "[{\"name\":\"newrelic-metrics-connector-hub-us-ashburn-1\", ...}]"
```

## Algorithm

Compartments are sorted by OCID, then packed using next-fit-with-splitting:

- Fill the active hub until it reaches the 5-compartment cap or the 50-namespace cap.
- When a compartment's namespaces would overflow the 50-namespace cap, split that compartment across two hubs.
- Close the active hub and open a new one when either cap is hit.

Sorting by OCID makes the output **deterministic**: re-running the script with the same input produces the same hub layout, so a subsequent `terraform apply` is a no-op.

## Hub naming

| Hubs produced | Names |
|---|---|
| 1 | `newrelic-metrics-connector-hub-{region}` |
| 2+ | `newrelic-metrics-connector-hub-{region}`, `newrelic-metrics-connector-hub-{region}-1`, ... |

## Example

Input (`compartments.json`):

```json
{
  "region": "us-ashburn-1",
  "compartments": [
    {
      "compartment_id": "ocid1.compartment.oc1..aaaa",
      "namespaces": ["oci_compute", "oci_lbaas"]
    },
    {
      "compartment_id": "ocid1.compartment.oc1..bbbb",
      "namespaces": ["oci_faas", "oci_blockstore", "oci_vcn"]
    }
  ]
}
```

Run:

```bash
python3 connector_hub_calculator.py --input compartments.json --pretty
```

Output:

```json
[
  {
    "name": "newrelic-metrics-connector-hub-us-ashburn-1",
    "description": "[DO NOT DELETE] New Relic Metrics Connector Hub",
    "compartments": [
      {
        "compartment_id": "ocid1.compartment.oc1..aaaa",
        "namespaces": ["oci_compute", "oci_lbaas"]
      },
      {
        "compartment_id": "ocid1.compartment.oc1..bbbb",
        "namespaces": ["oci_faas", "oci_blockstore", "oci_vcn"]
      }
    ]
  }
]
```
