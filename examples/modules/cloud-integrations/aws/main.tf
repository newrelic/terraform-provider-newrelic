data "aws_iam_policy_document" "newrelic_assume_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "AWS"
      // This is the unique identifier for New Relic account on AWS, there is no need to change this
      identifiers = [754728514883]
    }

    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"
      values   = [var.newrelic_account_id]
    }
  }
}

resource "aws_iam_role" "newrelic_aws_role" {
  name               = "NewRelicInfrastructure-Integrations-${var.name}"
  description        = "New Relic Cloud integration role"
  assume_role_policy = data.aws_iam_policy_document.newrelic_assume_policy.json
}

resource "aws_iam_policy" "newrelic_aws_permissions" {
  name        = "NewRelicCloudStreamReadPermissions-${var.name}"
  description = ""
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
        "ec2:DescribeInternetGateways",
        "ec2:DescribeVpcs",
        "ec2:DescribeNatGateways",
        "ec2:DescribeVpcEndpoints",
        "ec2:DescribeSubnets",
        "ec2:DescribeNetworkAcls",
        "ec2:DescribeVpcAttribute",
        "ec2:DescribeRouteTables",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVpcPeeringConnections",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DescribeVpnConnections",
        "health:DescribeAffectedEntities",
        "health:DescribeEventDetails",
        "health:DescribeEvents",
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

resource "aws_iam_role_policy_attachment" "newrelic_aws_policy_attach" {
  role       = aws_iam_role.newrelic_aws_role.name
  policy_arn = aws_iam_policy.newrelic_aws_permissions.arn
}

resource "aws_iam_role_policy_attachment" "newrelic_aws_policy_attach_read_only" {
  role       = alks_iamrole.newrelic_aws_role.name
  policy_arn = "arn:aws:iam::aws:policy/ReadOnlyAccess"
}

resource "newrelic_cloud_aws_link_account" "newrelic_cloud_integration_push" {
  account_id             = var.newrelic_account_id
  arn                    = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PUSH"
  name                   = "${var.name} metric stream"
  depends_on             = [aws_iam_role_policy_attachment.newrelic_aws_policy_attach]
}

resource "newrelic_api_access_key" "newrelic_aws_access_key" {
  account_id  = var.newrelic_account_id
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "Metric Stream Key for ${var.name}"
  notes       = "AWS Cloud Integrations Metric Stream Key"
}

resource "aws_iam_role" "firehose_newrelic_role" {
  name = "firehose_newrelic_role_${var.name}"

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
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket" "newrelic_aws_bucket" {
  bucket        = "newrelic-aws-bucket-${random_string.s3-bucket-name.id}"
  force_destroy = true
}

resource "aws_s3_bucket_ownership_controls" "newrelic_ownership_controls" {
  bucket = aws_s3_bucket.newrelic_aws_bucket.id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

locals {
  newrelic_urls = {
    US      = "https://aws-api.newrelic.com/cloudwatch-metrics/v1"
    EU      = "https://aws-api.eu01.nr-data.net/cloudwatch-metrics/v1"
  }
}

resource "aws_kinesis_firehose_delivery_stream" "newrelic_firehose_stream" {
  name        = "newrelic_firehose_stream_${var.name}"
  destination = "http_endpoint"
  http_endpoint_configuration {
    url                = local.newrelic_urls[var.newrelic_account_region]
    name               = "New Relic ${var.name}"
    access_key         = newrelic_api_access_key.newrelic_aws_access_key.key
    buffering_size     = 1
    buffering_interval = 60
    role_arn           = aws_iam_role.firehose_newrelic_role.arn
    s3_backup_mode     = "FailedDataOnly"
    s3_configuration {
      role_arn           = aws_iam_role.firehose_newrelic_role.arn
      bucket_arn         = aws_s3_bucket.newrelic_aws_bucket.arn
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
  name = "newrelic_metric_stream_to_firehose_role_${var.name}"

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
  name = "default"
  role = aws_iam_role.metric_stream_to_firehose.id

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
            "Resource": "${aws_kinesis_firehose_delivery_stream.newrelic_firehose_stream.arn}"
        }
    ]
}
EOF
}

resource "aws_cloudwatch_metric_stream" "newrelic_metric_stream" {
  name          = "newrelic-metric-stream-${var.name}"
  role_arn      = aws_iam_role.metric_stream_to_firehose.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.newrelic_firehose_stream.arn
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

resource "newrelic_cloud_aws_link_account" "newrelic_cloud_integration_pull" {
  account_id             = var.newrelic_account_id
  arn                    = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PULL"
  name                   = "${var.name} pull"
  depends_on             = [aws_iam_role_policy_attachment.newrelic_aws_policy_attach]
}

resource "newrelic_cloud_aws_integrations" "newrelic_cloud_integration_pull" {
  account_id        = var.newrelic_account_id
  linked_account_id = newrelic_cloud_aws_link_account.newrelic_cloud_integration_pull.id
  billing {}
  cloudtrail {}
  health {}
  trusted_advisor {}
  vpc {}
  x_ray {}
  s3 {}
  doc_db {}
  sqs {}
  ebs {}
  alb {}
  elasticache {}
  api_gateway {}
  auto_scaling {}
  aws_app_sync {}
  aws_athena {}
  aws_cognito {}
  aws_connect {}
  aws_direct_connect {}
  aws_fsx {}
  aws_glue {}
  aws_kinesis_analytics {}
  aws_media_convert {}
  aws_media_package_vod {}
  aws_mq {}
  aws_msk {}
  aws_neptune {}
  aws_qldb {}
  aws_route53resolver {}
  aws_states {}
  aws_transit_gateway {}
  aws_waf {}
  aws_wafv2 {}
  cloudfront {}
  dynamodb {}
  ec2 {}
  ecs {}
  efs {}
  elasticbeanstalk {}
  elasticsearch {}
  elb {}
  emr {}
  iam {}
  iot {}
  kinesis {}
  kinesis_firehose {}
  lambda {}
  rds {}
  redshift {}
  route53 {}
  ses {}
  sns {}
}

resource "aws_s3_bucket" "newrelic_configuration_recorder_s3" {
  bucket        = "newrelic-configuration-recorder-${random_string.s3-bucket-name.id}"
  force_destroy = true
}

locals {
  should_create_recorder = var.enable_config_recorder
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
        "${aws_s3_bucket.newrelic_configuration_recorder_s3.arn}",
        "${aws_s3_bucket.newrelic_configuration_recorder_s3.arn}/*"
      ]
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "newrelic_configuration_recorder" {
  count      = local.should_create_recorder ? 1 : 0
  role       = aws_iam_role.newrelic_configuration_recorder[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWS_ConfigRole"
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
  s3_bucket_name = aws_s3_bucket.newrelic_configuration_recorder_s3.bucket
  depends_on = [
    aws_config_configuration_recorder.newrelic_recorder
  ]
}
