---
layout: "newrelic"
page_title: "New Relic: newrelic_streaming_export_rule"
sidebar_current: "docs-newrelic-resource-streaming-export-rule"
description: |-
  Create and manage a streaming export rule in New Relic.
---

# Resource: newrelic\_streaming\_export\_rule

Use this resource to create and manage New Relic streaming export rules. Details regarding streaming export and permissions can be found [here](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-streaming-export/).

## Example Usage

##### AWS Kinesis Firehose
```hcl
resource "newrelic_streaming_export_rule" "foo" {
  account_id = 12345678
  name       = "foo"
  nrql       = "SELECT * FROM Transaction"
  description = "Example streaming export rule"

  payload_compression = "GZIP"

  aws {
    aws_account_id      = 123456789012
    delivery_stream_name = "my-delivery-stream"
    region              = "us-east-1"
    role                = "arn:aws:iam::123456789012:role/NewRelicStreamingExportRole"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account ID of the account that owns the streaming export rule. Defaults to the account associated with the API key used.
* `name` - (Required) The name of the streaming export rule.
* `nrql` - (Required) NRQL to select the telemetry data to export.
* `description` - (Optional) Additional information about the streaming export rule.
* `payload_compression` - (Optional) Whether to compress payloads before sending them out. One of: `DISABLED`, `GZIP`.
* `aws` - (Optional) A nested block that describes AWS parameters for the streaming export rule. Only one `aws` block is permitted per streaming export rule definition. See [Nested aws blocks](#nested-aws-blocks) below for details.
* `azure` - (Optional) A nested block that describes Azure parameters for the streaming export rule. Only one `azure` block is permitted per streaming export rule definition. See [Nested azure blocks](#nested-azure-blocks) below for details.
* `gcp` - (Optional) A nested block that describes GCP parameters for the streaming export rule. Only one `gcp` block is permitted per streaming export rule definition. See [Nested gcp blocks](#nested-gcp-blocks) below for details.

### Nested `aws` blocks

* `aws_account_id` - (Required) The AWS account to which the target firehose belongs.
* `delivery_stream_name` - (Required) The name of the delivery stream to write events to.
* `region` - (Required) The AWS region the delivery stream is located in.
* `role` - (Required) The role configured for New Relic to assume.

### Nested `azure` blocks

* `event_hub_connection_string` - (Required) Connection string that has access to the specific Event Hub.
* `event_hub_name` - (Required) The name of Event Hub to write events to.

### Nested `gcp` blocks

* `gcp_project_id` - (Required) The GCP project ID.
* `pubsub_topic_id` - (Required) The Pub/Sub topic ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the streaming export rule.
* `status` - The state of the streaming export rule. One of: `CREATION_FAILED`, `CREATION_IN_PROGRESS`, `DELETED`, `DISABLED`, `ENABLED`.
* `message` - A message returned by the latest API call.
* `created_at` - The time at which the process of creating the streaming rule began.
* `updated_at` - The last time the status of the streaming rule was updated.

## Additional Examples

~> **NOTE:** Exactly one of `aws`, `azure`, or `gcp` should be specified per streaming export rule.

#### Azure Event Hub

```hcl
resource "newrelic_streaming_export_rule" "foo" {
  account_id  = 12345678
  name        = "azure-streaming-export"
  nrql        = "SELECT * FROM Log"
  description = "Export logs to Azure Event Hub"

  azure {
    event_hub_connection_string = "Endpoint=sb://my-namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=abc123"
    event_hub_name              = "my-event-hub"
  }
}
```

#### GCP Pub/Sub

```hcl
resource "newrelic_streaming_export_rule" "foo" {
  account_id  = 12345678
  name        = "gcp-streaming-export"
  nrql        = "SELECT * FROM Metric"
  description = "Export metrics to GCP Pub/Sub"

  gcp {
    gcp_project_id = "my-gcp-project"
    pubsub_topic_id = "my-pubsub-topic"
  }
}
```

## Additional Information

More information about streaming export can be found in the New Relic [documentation](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-streaming-export/).