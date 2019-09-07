// The name of an existing notification channel you'd like to use
/*
data "newrelic_alert_channel" "team_alert_channel" {
  // MODIFY THIS LINE BELOW
  name = "ChangeMe alert channel"
}

# Configure the New Relic provider
data "aws_kms_secrets" "newrelic_api_key" {
  secret {
    name    = "api_key"
    // MODIFY THIS LINE BELOW - insert KMS payload from Cloud Services
    payload = ""
  }
}

terraform {
  backend "s3" {
    bucket = "nr-relfit-terraform-state"

    // MODIFY THIS LINE BELOW - Should be unique across each of your environments
    key      = "ChangeMe-production-state"
    region   = "us-east-1"
    role_arn = "ChangeMe state bucket arn from Cloud Services"
  }
}
*/
provider "newrelic" {
  api_key = "9cec0ab87a3a9b02ab10d910f3779736"
  // REVIEW THIS LINE IN CASE YOUR SERVICE NEEDS A DIFFERENT VALUE
  api_url = "https://api.newrelic.com/v2"
}
resource "newrelic_synthetics_monitor" "synthetics_monitor" {
  type      = "SIMPLE"
  name      = "terraform_example"
  frequency = 10
  uri       = "http://www.google.com"
  locations = ["AWS_US_WEST_1"]
  status    = "ENABLED"

}
