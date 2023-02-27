#Goal :

terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
      #   version = "3.0.3"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "4.54.0"
    }
  }
}

provider "newrelic" {
  # Configuration options
  region = "US"
}
provider "aws" {
  region = "us-west-2"

}


/*

    Complete example to enable AWS integration with New Relic

*/

variable "NEW_RELIC_ACCOUNT_ID" {
  type = string
  default = "2520528"
}

variable "NEW_RELIC_CLOUDWATCH_ENDPOINT" {
  type = string
  # Depending on where your New Relic Account is located you need to change the default
  default = "https://aws-api.newrelic.com/cloudwatch-metrics/v1" # US Datacenter
  # default = "https://aws-api.eu01.nr-data.net/cloudwatch-metrics/v1" # EU Datacenter
}

variable "NEW_RELIC_ACCOUNT_NAME" {
  type    = string
  default = "Production"
}

 data "aws_iam_policy_document" "newrelic_assume_policy" {
   statement {
     actions = ["sts:AssumeRole"]

     principals {
       type        = "AWS"
       // This is the unique identifier for New Relic account on AWS, there is no need to change this
       identifiers = [754728514883]
     }

     condition {
       test     = "StringEquals"
       variable = "sts:ExternalId"
       values   = [var.NEW_RELIC_ACCOUNT_ID]
     }
   }
 }


# removing it because it exists reading from above
 resource "aws_iam_role" "newrelic_aws_role" {
   name               = "NewRelicInfrastructure-Integrations"
   description        = "New Relic Cloud integration role"
   assume_role_policy = aws_iam_policy_document.newrelic_assume_policy.json
 }

 resource "aws_iam_policy" "newrelic_aws_permissions" {
   name        = "NewRelicCloudStreamReadPermissions"
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

 # 1 Iam Role = attaching a new policy with required config
 resource "aws_iam_role_policy_attachment" "newrelic_aws_policy_attach" {
   role       = aws_iam_role.newrelic_aws_role.name
   policy_arn = aws_iam_policy.newrelic_aws_permissions.arn
 }

# using 1 Iam role to link acc to 1 Nr acc for PUSH
resource "newrelic_cloud_aws_link_account" "newrelic_cloud_integration_push" {
  arn = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PUSH"
  name                   = "${var.NEW_RELIC_ACCOUNT_NAME} Push"
  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_policy_attach]
}

# new link to the same account for PULL same role - different name
resource "newrelic_cloud_aws_link_account" "newrelic_cloud_integration_pull" {
   account_id = var.NEW_RELIC_ACCOUNT_ID
  arn                    = aws_iam_role.newrelic_aws_role.arn
  metric_collection_mode = "PULL"
  name                   = "${var.NEW_RELIC_ACCOUNT_NAME} Pull"
  depends_on = [aws_iam_role_policy_attachment.newrelic_aws_policy_attach]
}

# Get the API key
resource "newrelic_api_access_key" "newrelic_aws_access_key" {
  account_id  = var.NEW_RELIC_ACCOUNT_ID
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "Ingest License key"
  notes       = "AWS Cloud Integrations Firehost Key"
}

# add firehose service to role
 resource "aws_iam_role" "firehose_newrelic_role" {
   name = "firehose_newrelic_role"

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

# #RANDOM STRING gen
resource "random_string" "s3-bucket-name" {
  length  = 8
  special = false
  upper   = false
}

# creating a s3 bucket
resource "aws_s3_bucket" "newrelic_aws_bucket" {
  bucket = "newrelic-aws-bucket-${random_string.s3-bucket-name.id}"
}


resource "aws_s3_bucket_acl" "newrelic_aws_bucket_acl" {
  bucket = aws_s3_bucket.newrelic_aws_bucket.id
  acl    = "private"
}

 # firehose to send data to newrelic account
 resource "aws_kinesis_firehose_delivery_stream" "newrelic_firehost_stream" {
   name        = "newrelic_firehost_stream"
   destination = "http_endpoint"

   s3_configuration {
     role_arn           = aws_iam_role.firehose_newrelic_role.arn
     bucket_arn         = aws_s3_bucket.newrelic_aws_bucket.arn
     buffer_size        = 10
     buffer_interval    = 400
     compression_format = "GZIP"
   }

   http_endpoint_configuration {
     url                = var.NEW_RELIC_CLOUDWATCH_ENDPOINT
     name               = "New Relic"
     access_key         = newrelic_api_access_key.newrelic_aws_access_key.key
     buffering_size     = 1
     buffering_interval = 60
     role_arn           = aws_iam_role.firehose_newrelic_role.arn
     s3_backup_mode     = "FailedDataOnly"

     request_configuration {
       content_encoding = "GZIP"
     }
   }
 }

data "aws_kinesis_firehose_delivery_stream" "newrelic_firehost_stream" {
  name = "newrelic_firehost_stream"
}

# # adding cloudwatch service to role
 resource "aws_iam_role" "metric_stream_to_firehose" {
   name = "metric_stream_to_firehose_role"

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

# ? sending/putting data to firehose
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
            "Resource": "${data.aws_kinesis_firehose_delivery_stream.newrelic_firehost_stream.arn}"
        }
    ]
}
EOF
}

# creating metric stream
resource "aws_cloudwatch_metric_stream" "newrelic_metric_stream" {
  name          = "newrelic-metric-stream"
  role_arn      = aws_iam_role.metric_stream_to_firehose.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.newrelic_firehost_stream.arn
  output_format = "opentelemetry0.7"
}

# integrate aws to nr
resource "newrelic_cloud_aws_integrations" "foo" {
  account_id        = var.NEW_RELIC_ACCOUNT_ID
  linked_account_id = newrelic_cloud_aws_link_account.newrelic_cloud_integration_pull.id
  billing {}
  cloudtrail {}
  health {}
  trusted_advisor {}
  vpc {}
  x_ray {}
  # aws_rds{}
}

# creating a aws rds
resource "aws_db_instance" "default" {
  allocated_storage    = 10
  db_name              = "TfTest${var.NEW_RELIC_ACCOUNT_ID}"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true
}

# # create a new s3 bucket
resource "aws_s3_bucket" "newrelic_configuration_recorder_s3" {
  bucket = "newrelic-configuration-recorder-${random_string.s3-bucket-name.id}"
}

# add new role to it with aws config
# AWS Config provides a detailed view of the configuration of AWS resources in your AWS account.
resource "aws_iam_role" "newrelic_configuration_recorder" {
  name = "newrelic_configuration_recorder-${var.NEW_RELIC_ACCOUNT_NAME}"

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

# aws iam role for s3 config
# Grant AWS Config the permissions it needs to access the Amazon S3 bucket and the Amazon SNS topic.
resource "aws_iam_role_policy" "newrelic_configuration_recorder_s3" {
  name = "newrelic-configuration-recorder-s3-${var.NEW_RELIC_ACCOUNT_NAME}"
  role = aws_iam_role.newrelic_configuration_recorder.id

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

# attaching the new role to an existing policy
resource "aws_iam_role_policy_attachment" "newrelic_configuration_recorder" {
  role       = aws_iam_role.newrelic_configuration_recorder.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWS_ConfigRole"
}

#AWS Config uses the configuration recorder to detect changes in your resource configurations
#and capture these changes as configuration items.
#You must create a configuration recorder before AWS Config can track your resource configurations.
resource "aws_config_configuration_recorder" "newrelic_recorder" {
  name     = "newrelic_configuration_recorder-${var.NEW_RELIC_ACCOUNT_NAME}"
  role_arn = aws_iam_role.newrelic_configuration_recorder.arn
}

# setting status of the config recorder
resource "aws_config_configuration_recorder_status" "newrelic_recorder_status" {
  name       = aws_config_configuration_recorder.newrelic_recorder.name
  is_enabled = true
  depends_on = [aws_config_delivery_channel.newrelic_recorder_delivery]
}

# Set up an Amazon S3 bucket to receive a configuration snapshot on request and configuration history.
resource "aws_config_delivery_channel" "newrelic_recorder_delivery" {
  name           = "newrelic_configuration_recorder-${var.NEW_RELIC_ACCOUNT_NAME}"
  s3_bucket_name = aws_s3_bucket.newrelic_configuration_recorder_s3.bucket
  depends_on = [
    aws_config_configuration_recorder.newrelic_recorder
  ]
}
