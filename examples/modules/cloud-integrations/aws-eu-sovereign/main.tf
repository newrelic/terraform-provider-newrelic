# AWS EU Sovereign Cloud Integration Module
# This module creates a New Relic AWS EU Sovereign cloud integration
# EU Sovereign only supports: billing, cloudtrail, health, trusted_advisor, x_ray

locals {
  # EU Sovereign only supports EU region
  newrelic_metric_api_url = "https://aws-api.eu01.nr-data.net/cloudwatch-metrics/v1"
  create_push             = contains(["PUSH", "BOTH"], var.metric_collection_mode)
  create_pull             = contains(["PULL", "BOTH"], var.metric_collection_mode)
  should_create_recorder  = var.enable_config_recorder
}

# IAM policy document for New Relic to assume the role
data "aws_iam_policy_document" "newrelic_assume_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "AWS"
      # This is the unique identifier for New Relic account on AWS EU Sovereign Cloud
      identifiers = ["093926027544"]
    }

    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"
      values   = [var.newrelic_account_id]
    }
  }
}

# Create IAM role for New Relic
resource "aws_iam_role" "newrelic_aws_role" {
  name               = "NewRelicInfrastructure-Integrations-${var.name}"
  description        = "New Relic Cloud integration role for EU Sovereign"
  assume_role_policy = data.aws_iam_policy_document.newrelic_assume_policy.json
}

# IAM policy with permissions for New Relic integrations
resource "aws_iam_policy" "newrelic_aws_permissions" {
  name        = "NewRelicCloudStreamReadPermissions-${var.name}"
  description = "Permissions for New Relic AWS EU Sovereign Cloud integrations"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "budgets:ViewBudget",
        "cloudtrail:LookupEvents",
        "config:BatchGetResourceConfig",
        "config:ListDiscoveredResources",
        "health:DescribeAffectedEntities",
        "health:DescribeEventDetails",
        "health:DescribeEvents",
        "support:DescribeTrustedAdvisorCheckRefreshStatuses",
        "support:DescribeTrustedAdvisorCheckResult",
        "support:DescribeTrustedAdvisorCheckSummaries",
        "support:DescribeTrustedAdvisorChecks",
        "support:RefreshTrustedAdvisorCheck",
        "tag:GetResources",
        "xray:BatchGet*",
        "xray:Get*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# Attach the policy to the role
resource "aws_iam_role_policy_attachment" "newrelic_aws_policy_attach" {
  role       = aws_iam_role.newrelic_aws_role.name
  policy_arn = aws_iam_policy.newrelic_aws_permissions.arn
}

# Wait for IAM role to propagate in EU Sovereign region
# EU Sovereign has slower IAM propagation than standard AWS/GovCloud
resource "terraform_data" "wait_for_iam" {
  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_policy_attach]

  provisioner "local-exec" {
    command = "sleep 10"
  }
}

# PUSH mode resources (Metric Streams)
resource "newrelic_cloud_aws_eu_sovereign_link_account" "newrelic_cloud_integration_push" {
  count                  = local.create_push ? 1 : 0
  account_id             = var.newrelic_account_id
  arn                    = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PUSH"
  name                   = "${var.name} metric stream"
  depends_on             = [terraform_data.wait_for_iam]
}

resource "newrelic_api_access_key" "newrelic_aws_access_key" {
  count       = local.create_push ? 1 : 0
  account_id  = var.newrelic_account_id
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "Metric Stream Key for ${var.name}"
  notes       = "AWS EU Sovereign Cloud Integrations Metric Stream Key"
}

resource "aws_iam_role" "firehose_newrelic_role" {
  count = local.create_push ? 1 : 0
  name  = "firehose_newrelic_role_${var.name}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "random_string" "s3-bucket-name" {
  count   = local.create_push ? 1 : 0
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket" "newrelic_aws_bucket" {
  count         = local.create_push ? 1 : 0
  bucket        = "newrelic-aws-bucket-${random_string.s3-bucket-name[0].id}"
  force_destroy = true
}

resource "aws_s3_bucket_ownership_controls" "newrelic_ownership_controls" {
  count  = local.create_push ? 1 : 0
  bucket = aws_s3_bucket.newrelic_aws_bucket[0].id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_kinesis_firehose_delivery_stream" "newrelic_firehose_stream" {
  count       = local.create_push ? 1 : 0
  name        = "newrelic_firehose_stream_${var.name}"
  destination = "http_endpoint"
  http_endpoint_configuration {
    url                = local.newrelic_metric_api_url
    name               = "New Relic ${var.name}"
    access_key         = newrelic_api_access_key.newrelic_aws_access_key[0].key
    buffering_size     = 1
    buffering_interval = 60
    role_arn           = aws_iam_role.firehose_newrelic_role[0].arn
    s3_backup_mode     = "FailedDataOnly"
    s3_configuration {
      role_arn           = aws_iam_role.firehose_newrelic_role[0].arn
      bucket_arn         = aws_s3_bucket.newrelic_aws_bucket[0].arn
      buffering_size     = 10
      buffering_interval = 400
      compression_format = "GZIP"
    }
    request_configuration {
      content_encoding = "GZIP"
    }
  }
}

resource "aws_iam_role" "metric_stream_to_firehose" {
  count = local.create_push ? 1 : 0
  name  = "newrelic_metric_stream_to_firehose_role_${var.name}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "streams.metrics.cloudwatch.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "metric_stream_to_firehose" {
  count = local.create_push ? 1 : 0
  name  = "default"
  role  = aws_iam_role.metric_stream_to_firehose[0].id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "firehose:PutRecord",
                "firehose:PutRecordBatch"
            ],
            "Resource": "${aws_kinesis_firehose_delivery_stream.newrelic_firehose_stream[0].arn}"
        }
    ]
}
EOF
}

resource "aws_cloudwatch_metric_stream" "newrelic_metric_stream" {
  count         = local.create_push ? 1 : 0
  name          = "newrelic-metric-stream-${var.name}"
  role_arn      = aws_iam_role.metric_stream_to_firehose[0].arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.newrelic_firehose_stream[0].arn
  output_format = var.output_format

  dynamic "exclude_filter" {
    for_each = var.exclude_metric_filters

    content {
      namespace    = exclude_filter.key
      metric_names = exclude_filter.value
    }
  }

  dynamic "include_filter" {
    for_each = var.include_metric_filters

    content {
      namespace    = include_filter.key
      metric_names = include_filter.value
    }
  }
}

# PULL mode resources (API Polling)
resource "newrelic_cloud_aws_eu_sovereign_link_account" "newrelic_cloud_integration_pull" {
  count                  = local.create_pull ? 1 : 0
  account_id             = var.newrelic_account_id
  arn                    = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PULL"
  name                   = "${var.name} pull"
  depends_on             = [terraform_data.wait_for_iam]
}

# Configure AWS EU Sovereign integrations
# EU Sovereign only supports: billing, cloudtrail, health, trusted_advisor, and x_ray
resource "newrelic_cloud_aws_eu_sovereign_integrations" "newrelic_cloud_integration_pull" {
  count             = local.create_pull ? 1 : 0
  account_id        = var.newrelic_account_id
  linked_account_id = newrelic_cloud_aws_eu_sovereign_link_account.newrelic_cloud_integration_pull[0].id
  billing {}
  cloudtrail {}
  health {}
  trusted_advisor {}
  x_ray {}
}

# Configuration Recorder resources
resource "random_string" "config_recorder_bucket_name" {
  count   = local.should_create_recorder ? 1 : 0
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket" "newrelic_configuration_recorder_s3" {
  count         = local.should_create_recorder ? 1 : 0
  bucket        = "newrelic-configuration-recorder-${random_string.config_recorder_bucket_name[0].id}"
  force_destroy = true
}

resource "aws_iam_role" "newrelic_configuration_recorder" {
  count = local.should_create_recorder ? 1 : 0
  name  = "newrelic_configuration_recorder-${var.name}"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "config.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
      ]
    }
EOF
}

resource "aws_iam_role_policy" "newrelic_configuration_recorder_s3" {
  count = local.should_create_recorder ? 1 : 0
  name  = "newrelic-configuration-recorder-s3-${var.name}"
  role  = aws_iam_role.newrelic_configuration_recorder[0].id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_s3_bucket.newrelic_configuration_recorder_s3[0].arn}",
        "${aws_s3_bucket.newrelic_configuration_recorder_s3[0].arn}/*"
      ]
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "newrelic_configuration_recorder" {
  count      = local.should_create_recorder ? 1 : 0
  role       = aws_iam_role.newrelic_configuration_recorder[0].name
  policy_arn = "arn:aws-eusc:iam::aws:policy/service-role/AWS_ConfigRole"
}

resource "aws_config_configuration_recorder" "newrelic_recorder" {
  count    = local.should_create_recorder ? 1 : 0
  name     = "newrelic_configuration_recorder-${var.name}"
  role_arn = aws_iam_role.newrelic_configuration_recorder[0].arn
}

resource "aws_config_configuration_recorder_status" "newrelic_recorder_status" {
  count      = local.should_create_recorder ? 1 : 0
  name       = aws_config_configuration_recorder.newrelic_recorder[0].name
  is_enabled = true
  depends_on = [aws_config_delivery_channel.newrelic_recorder_delivery]
}

resource "aws_config_delivery_channel" "newrelic_recorder_delivery" {
  count          = local.should_create_recorder ? 1 : 0
  name           = "newrelic_configuration_recorder-${var.name}"
  s3_bucket_name = aws_s3_bucket.newrelic_configuration_recorder_s3[0].bucket
  depends_on = [
    aws_config_configuration_recorder.newrelic_recorder
  ]
}