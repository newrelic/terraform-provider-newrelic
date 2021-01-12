---
layout: "newrelic"
page_title: "New Relic Provider Configuration"
sidebar_current: "docs-newrelic-provider-configuration"
description: |-
  Details on how to configure the New Relic Terraform Provider.
---

# Configuring the New Relic Terraform Provider

The New Relic Terraform Provider supports two methods of configuring the provider.

1. [Using environment variables](#configuration-via-environment-variables) (recommended)
1. [Using the `provider` block](#configuration-via-the-provider-block)

-> <small>**Note:** If you need to configure more than one instance of the New Relic provider, such as for different regions, we've provided an [example](#configuring-multiple-instances-of-the-provider) showing how this can be accomplished.</small>

### Configuration via environment variables

Certain [environment variables](#environment-variables-reference) will be automatically detected by the New Relic Terraform Provider when running `terraform` commands. Using environment variables facilitates setting a default provider configuration, resulting in a smaller `provider` block in your HCL, and it also keeps your credentials out of your source code.

If you're using Terraform locally, you can set the environment variables in your machine's startup file, such as your `.bash_profile` or `.bashrc` file on UNIX machines.

**.bash_profile**

```bash
# Add this to your .bash_profile or .bashrc
export NEW_RELIC_API_KEY="<your New Relic User API key>"
export NEW_RELIC_REGION="US"
```

Provided you have the required environment variables set, your `provider` block can look as simple as the following.

```hcl
provider "newrelic" {}
```

-> <small>**Note:** When using Terraform in your CI/CD pipeline, we recommend setting your environment variables within your platform's secrets management. Each platform, such as GitHub or CircleCI, has their own way of managing secrets and environment variables, so you will need to refer to your vendor's documentation for implemenation details.</small>


#### Environment variables reference

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

-> <small>**Note:** The `provider` block schema attributes take precedence over environment variables, providing the ability to override environment variables if needed. This can useful when using [multiple instances of the provider](#configuring-multiple-instances-of-the-provider).</small>


### Configuration via the `provider` block

Configuring the provider from within your HCL is a quick way to get started, however, we recommend [using environment variables](#configuration-via-environment-variables). The minimal recommended configuration is as follows.

```hcl
provider "newrelic" {
  api_key = <Your User API key>
  account_id = <your New Relic account ID>
  region = "US"
}
```


### Configuring multiple instances of the provider

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

-> <small>**Note:** The provider supports ***one*** region per instance of the provider.</small>

[account ID]: https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/account-id
[User API key]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#user-api-key
[data center region]: https://docs.newrelic.com/docs/using-new-relic/welcome-new-relic/get-started/our-eu-us-region-data-centers
[Insights query API key]: https://docs.newrelic.com/docs/insights/insights-api/get-data/query-insights-event-data-api
[Insights insert API key]: https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/introduction-event-api#register
