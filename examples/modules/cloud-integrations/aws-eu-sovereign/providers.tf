terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.22.0"  # EU Sovereign requires <= 6.22.0 due to S3 tagging API bug in 6.23.0+
    }
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}