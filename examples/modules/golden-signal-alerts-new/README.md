# [Golden Signal Alerts (New)](modules/golden-signal-alerts-new)
This module encapsulates an alerting strategy based on the [Four Golden Signals](https://landing.google.com/sre/sre-book/chapters/monitoring-distributed-systems/#xref_monitoring_golden-signals) introduced in Google’s widely read book on [Site Reliability Engineering](https://landing.google.com/sre/sre-book/toc/index.html).

The signals chosen for this module are:

* *Latency*: High response time (seconds)
* *Traffic*: Low throughput (requests/minute)
* *Errors*: Error rate (errors/minute)
* *Saturation*: CPU utilization (percentage utilized)

### Requirements
Applications making use of this module need to be reporting data into both APM and Infrastructure.

### Input variables
The following input variables are accepted by the module:

* `name`: The APM application name as reported to New Relic
* `threshold_duration`: The duration, in seconds, that the condition must violate the threshold before creating a violation.
* `cpu_threshold`: The critical threshold of the CPU utilization condition, as a percentage
* `error_percentage_threshold`: The critical threshold of the error rate condition, as a percentage
* `response_time_threshold`: The critical threshold of the response time condition, in seconds
* `throughput_threshold`: The critical threshold of the throughput condition, in requests/second

### Outputs
The following output values are provided by the module:

* `policy_id`: The ID of the created alert policy
* `cpu_condition_id`: The ID of the created high CPU alert condition
* `error_percentage_condition_id`: The ID of the created error percentage alert condition
* `response_time_condition_id`: The ID of the created response time alert condition
* `throughput_condition_id`: The ID of the created throughput alert condition


### Example usage
```terraform

data "newrelic_notification_destination" "webhook_destination" {
  name = "Golden Signal Webhook Testing"
}

# Resource
resource "newrelic_notification_channel" "webhook_notification_channel" {
  name           = "webhook-example"
  type           = "WEBHOOK"
  destination_id = data.newrelic_notification_destination.webhook_destination.id
  product        = "IINT"

  property {
    key   = "payload"
    value = "{\n\t\"name\": \"foo\"\n}"
    label = "Payload Template"
  }
}

data "newrelic_notification_destination" "email_destination" {
  name = "golden signals testing mail"
}

resource "newrelic_notification_channel" "email_notification_channel" {
  name = "email-example"
  type = "EMAIL"
  destination_id = data.newrelic_notification_destination.email_destination.id
  product = "IINT"

  property {
    key = "subject"
    value = "New Subject Title"
  }

  property {
    key = "customDetailsEmail"
    value = "issue id - {{issueId}}"
  }
}

module "webportal_alerts" {
  // Please specify the path of the source of this module according to the location you've placed the module in.
  // The path specified below assumes you're using this module from a clone of this repo, in the `newrelic.tf` file in the `testing` folder.
  // However, if you'd like to use a remote version of this module (without a cloned version of this), the right value of the argument source would be "github.com/newrelic/terraform-provider-newrelic//examples/modules/golden-signal-alerts-new".
  source = "../examples/modules/golden-signal-alerts-new"
  notification_channel_ids = [newrelic_notification_channel.webhook_notification_channel.id, newrelic_notification_channel.email_notification_channel.id]

  service = {
    name                       = "Dummy App Pro Max"
    threshold_duration         = 420
    cpu_threshold              = 90
    response_time_threshold    = 5
    error_percentage_threshold = 10
    throughput_threshold       = 300
  }
}
```
