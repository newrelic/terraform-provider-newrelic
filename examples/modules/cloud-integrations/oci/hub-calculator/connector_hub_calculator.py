#!/usr/bin/env python3
"""
OCI Metrics Connector Hub Calculator

Computes the connector_hubs_data variable for the metrics-integration Terraform module.

OCI caps per Service Connector Hub monitoring source:
  - at most 5 compartments
  - at most 50 namespaces (summed across all compartments in the hub)

Uses next-fit-with-splitting over a compartment-OCID-sorted input, matching the
algorithm in the New Relic backend (beyond-api-v2 OciGenerateInstrumentedPayloadUrl).
Sorting makes the output deterministic: re-running with unchanged input produces
the same hub layout, so a subsequent terraform apply is a no-op.
"""

import argparse
import json
import sys

MAX_COMPARTMENTS_PER_HUB = 5
MAX_NAMESPACES_PER_HUB = 50
HUB_DESCRIPTION = "[DO NOT DELETE] New Relic Metrics Connector Hub"


def hub_name(region: str, index: int) -> str:
    base = f"newrelic-metrics-connector-hub-{region}"
    return base if index == 0 else f"{base}-{index}"


def pack(compartments: list[dict], region: str) -> list[dict]:
    """
    Pack compartments into connector hubs.

    Each item in `compartments` must have:
      compartment_id: str
      namespaces:     list[str]
    """
    sorted_compartments = sorted(compartments, key=lambda c: c["compartment_id"])

    hubs: list[dict] = []
    current_compartments: list[dict] = []
    current_ns_count = 0

    def close_hub():
        nonlocal current_compartments, current_ns_count
        if current_compartments:
            idx = len(hubs)
            hubs.append({
                "name": hub_name(region, idx),
                "description": HUB_DESCRIPTION,
                "compartments": current_compartments,
            })
        current_compartments = []
        current_ns_count = 0

    for comp in sorted_compartments:
        compartment_id = comp["compartment_id"]
        namespaces = list(comp["namespaces"])
        offset = 0
        remaining = len(namespaces)

        while remaining > 0:
            if len(current_compartments) >= MAX_COMPARTMENTS_PER_HUB:
                close_hub()

            slack = MAX_NAMESPACES_PER_HUB - current_ns_count

            if slack >= remaining:
                current_compartments.append({
                    "compartment_id": compartment_id,
                    "namespaces": namespaces[offset:offset + remaining],
                })
                current_ns_count += remaining
                remaining = 0
            elif slack > 0:
                current_compartments.append({
                    "compartment_id": compartment_id,
                    "namespaces": namespaces[offset:offset + slack],
                })
                close_hub()
                offset += slack
                remaining -= slack
            else:
                close_hub()

    close_hub()
    return hubs


def load_input(path: str) -> dict:
    with open(path) as f:
        return json.load(f)


def main():
    parser = argparse.ArgumentParser(
        description="Compute connector_hubs_data for the OCI metrics-integration Terraform module.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Input JSON format:
  {
    "region": "us-ashburn-1",
    "compartments": [
      {
        "compartment_id": "ocid1.compartment.oc1..aaa",
        "namespaces": ["oci_compute", "oci_lbaas"]
      }
    ]
  }

Output: a JSON string ready to paste into connector_hubs_data in your .tfvars file.

Example:
  python3 connector_hub_calculator.py --input compartments.json
  python3 connector_hub_calculator.py --input compartments.json --pretty
""",
    )
    parser.add_argument("--input", required=True, help="Path to input JSON file")
    parser.add_argument(
        "--pretty", action="store_true", help="Pretty-print the output instead of a single-line string"
    )
    args = parser.parse_args()

    try:
        data = load_input(args.input)
    except FileNotFoundError:
        print(f"error: input file not found: {args.input}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"error: invalid JSON in input file: {e}", file=sys.stderr)
        sys.exit(1)

    region = data.get("region")
    compartments = data.get("compartments", [])

    if not region:
        print("error: 'region' is required in input JSON", file=sys.stderr)
        sys.exit(1)
    if not compartments:
        print("error: 'compartments' list is empty", file=sys.stderr)
        sys.exit(1)

    for i, c in enumerate(compartments):
        if not c.get("compartment_id"):
            print(f"error: compartments[{i}] is missing 'compartment_id'", file=sys.stderr)
            sys.exit(1)
        if not c.get("namespaces"):
            print(f"error: compartments[{i}] has no namespaces", file=sys.stderr)
            sys.exit(1)

    hubs = pack(compartments, region)

    if args.pretty:
        print(json.dumps(hubs, indent=2))
    else:
        # Terraform expects connector_hubs_data as a JSON-encoded string
        print(json.dumps(json.dumps(hubs, separators=(",", ":"))))


if __name__ == "__main__":
    main()
