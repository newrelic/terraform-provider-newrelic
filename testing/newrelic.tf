terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

resource "newrelic_one_dashboard" "coaa_ghe" {
  name        = "NR-121612 GTSE Dashboard"
  permissions = "public_read_write"



  # ----------------------------------------------------------------------

  page {
    name = "GitHub Actions"


    widget_line {
      title  = "SQS - ages of oldest message"
      row    = 10
      column = 1
      height = 4
      width  = 6

      nrql_query {
        account_id = 3806526
        query      = <<EOT
FROM QueueSample SELECT latest(provider.approximateAgeOfOldestMessage.Maximum) as github_actions_runners_queued_builds
  WHERE provider = 'SqsQueue' AND provider.queueName = 'github-actions-runners-queued-builds.fifo'
  TIMESERIES MAX
EOT
      }
    }

    # ---------------------------------------------------------------------- right side charts


    widget_line {
      title  = "EC2 Instance Counts"
      row    = 7
      column = 7
      height = 4
      width  = 6

      y_axis_left_zero = true
      null_values { null_value = "remove" }
      nrql_query {
        account_id = 3806526
        query      = <<EOT
FROM Metric SELECT rate(uniqueCount(aws.ec2.InstanceId), 10 minutes) AS legacy_prod
  WHERE aws.accountId = '652911051897' AND aws.ec2.state = 'running' AND tags.Name LIKE '%github-actions%'
  TIMESERIES 10 minutes SLIDE BY MAX
EOT
      }
      nrql_query {
        account_id = 3806526
        query      = <<EOT
FROM Metric SELECT rate(uniqueCount(aws.ec2.InstanceId), 10 minutes) AS eks
  WHERE aws.accountId = '856289458132' AND aws.ec2.state = 'running' AND tags.Function = 'actions-runners'
  TIMESERIES 10 minutes SLIDE BY MAX
EOT
      }
    }



  }
}
