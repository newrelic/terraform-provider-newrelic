---
layout: "newrelic"
page_title: "New Relic Provider Configuration"
sidebar_current: "docs-newrelic-provider-configuration"
description: |-
  Details on how to configure the New Relic Terraform Provider.
---

# Configuring the New Relic Terraform Provider

The New Relic Terraform Provider supports two methods of configuring the provider.

1. [Using the `provider` block](#configuration-via-the-provider-block)
2. [Using environment variables](#configuration-via-environment-variables)

-> <small>If you need to configure more than one instance of the New Relic provider, such as for different regions, we've provided an [example](#configuring-multiple-instances-of-the-provider) showing how this can be accomplished.</small>

## Configuration via the `provider` block

Configuring the provider from within your HCL is a quick way to get started, however, we recommend [using environment variables](#configuration-via-environment-variables). The minimal recommended configuration is as follows.

```hcl
# Configure terraform
terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

# Configure the New Relic provider
provider "newrelic" {
  account_id = <Your Account ID>
  api_key = <Your User API Key>    # Usually prefixed with 'NRAK'
  region = "US"                    # Valid regions are US and EU
}
```

## Setting the New Relic Provider version

Terraform's syntax for [version constraints](https://www.terraform.io/language/expressions/version-constraints#version-constraint-syntax) is similar to the syntax used by other dependency management systems like Bundler and NPM.

In the [`provider` block configuration example](#configuration-via-the-provider-block), you will notice a version specification for the New Relic provider - i.e. `version = "~> 2.44"`. This will ensure Terraform always downloads and uses the latest *minor* version of the New Relic provider when running `terraform init`. If you prefer to only allow for patch updates, use `version = "~> 2.44.0"` (note the inclusion of the patch version digit).


## Upgrading the New Relic provider

1. Check for the [latest version of the New Relic provider](https://registry.terraform.io/providers/newrelic/newrelic/latest) in the Terraform registry.
2. Update your version in the `required_providers` configuration block. e.g. `"~> 2.0"`.
    ```hcl
    terraform {
      required_providers {
        newrelic = {
          source  = "newrelic/newrelic"
          version = "~> 2.44"  # Update this line
        }
      }
    }
    ```
    Using the configuration above as a hypothetical example, Terraform will download the latest `2.x` version of the provider, for example `2.45.0` if it exists. If `2.45.1` exists, it will download `2.45.1`. It will *not* download a `3.x` version.

    If a new major version of the provider exists and you need to upgrade, such as a version `3.0.0` or higher, the version string must be updated to `"~> 3.0"`, and from this point, Terraform will download the latest `3.x` version of the provider, for example `3.1.0`

3. Run `terraform init` to download and initialize the provider.

-> <small>Changing the version constraint string requires `terraform init` to be run to download and initialize the provider based on your updated version string.</small>


## Configuration via environment variables

Certain [environment variables](#environment-variables-reference) will be automatically detected by the New Relic Terraform Provider when running `terraform` commands. Using environment variables facilitates setting a default provider configuration, resulting in a smaller `provider` block in your HCL, and it also keeps your credentials out of your source code.

If you're using Terraform locally, you can set the environment variables in your machine's startup file, such as your `.bash_profile` or `.bashrc` file on UNIX machines.

**.bash_profile**

```bash
# Add this to your .bash_profile or .bashrc
export NEW_RELIC_API_KEY="<your New Relic User API key>"
export NEW_RELIC_REGION="US"
```

Provided you have the required environment variables set, your `provider` block can be left empty like the example below.

```hcl
provider "newrelic" {}
```

<br>

-> <small>When using Terraform in your CI/CD pipeline, we recommend setting your environment variables within your platform's secrets management. Each platform, such as GitHub or CircleCI, has their own way of managing secrets and environment variables, so you will need to refer to your vendor's documentation for implemenation details.</small>


## Environment variables reference

The table below shows the available environment variables and how they map to the provider's schema attributes. When using environment variables, you do *not* need to set the schema attributes within your `provider` block. All schema attributes default to their equivalent environment variables.

| <small>Schema Attribute</small> | <small>Equivalent Env Variable</small> | <small>Required?</small> | <small>Default</small> | <small>Description</small>                                                                   |
| ------------------------------- | -------------------------------------- | ------------------------ | ---------------------- | -------------------------------------------------------------------------------------------- |
| `account_id`                    | `NEW_RELIC_ACCOUNT_ID`                 | required                 | `null`                 | Your New Relic [account ID].                                                                 |
| `api_key`                       | `NEW_RELIC_API_KEY`                    | required                 | `null`                 | Your New Relic [User API key] \(usually prefixed with `NRAK`).                                     |
| `region`                        | `NEW_RELIC_REGION`                     | required                 | `null`                 | Your New Relic account's [data center region] \(`US` or `EU`).                               |
| `insights_insert_key`           | `NEW_RELIC_INSIGHTS_INSERT_KEY`        | optional                 | `null`                 | Your [Insights insert API key] for Insights events.                                          |
| `insecure_skip_verify`          | `NEW_RELIC_API_SKIP_VERIFY`            | optional                 | `null`                 | Whether or not to trust self-signed SSL certificates.                                        |
| `cacert_file`                   | `NEW_RELIC_API_CACERT`                 | optional                 | `null`                 | A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. |

<br>

-> <small>The `provider` block schema attributes take precedence over environment variables, providing the ability to override environment variables if needed. This can useful when using [multiple instances of the provider](#configuring-multiple-instances-of-the-provider).</small>


## Configuring multiple instances of the provider

The example below shows how you could use environment variables for the default configuration, then override the environment variables for another instance of the provider.

```hcl
# Default provider configuration via environment variables
provider "newrelic" {}

# Second provider for the EU region
provider "newrelic" {
  alias  = "europe"
  region = "EU"
}
```

-> <small>The provider supports ***one*** region per instance of the provider.</small>

[account ID]: https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/account-id
[User API key]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key
[data center region]: https://docs.newrelic.com/docs/using-new-relic/welcome-new-relic/get-started/our-eu-us-region-data-centers
[Insights query API key]: https://docs.newrelic.com/docs/insights/insights-api/get-data/query-insights-event-data-api
[Insights insert API key]: https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/introduction-event-api#register
