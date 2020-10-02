---
layout: "newrelic"
page_title: "New Relic Terraform Provider v2.x Migration Guide"
sidebar_current: "docs-newrelic-provider-v2-migration-guide"
description: |-
  Use this guide to update the New Relic Terraform Provider from v1.x to v2.x
---

## Upgrade to v2.x of the New Relic Terraform Provider

Version 2.0 of the provider introduces some changes to the provider's configuration. Users wanting to upgrade from v1.x to v2.x will need to make a few adjustments to their configuration prior to upgrading.

### A Note About API Key Format

Your New Relic [**Personal API Key**](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key) is now considered the default and standard API key for the provider.

New Relic [**Admin API keys**](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#admin) are no longer used for authentication within the New Rrelic Terraform provider.

-> <small>**Please note the following formatting for the provider's API key.** <br>Your **Personal API Key** has a prefix of `NRAK-` </small>

### Environment Variable Updates

If you have been using environment variables to configure the provider, you will need to take note of the following updates and make the necessary changes to your environment variables.

1. **IMPORTANT:** All environment variables in use by the provider have been renamed with a new naming convention. The `NEWRELIC_*` prefix has been changed to `NEW_RELIC_*`. This will be the naming convention for environment variables moving forward.

2. The environment variable `NEWRELIC_PERSONAL_API_KEY` has been replaced with `NEW_RELIC_API_KEY`. The Personal API Key is now considered the default and standard API key for the provider.

    ```diff
    - NEWRELIC_PERSONAL_API_KEY
    + NEW_RELIC_API_KEY
    ```

3. Use of `NEW_RELIC_ADMIN_API_KEY` has been discontinued.

    ```diff
    - NEW_RELIC_ADMIN_API_KEY
    ```

4. The v2 provider configuration **requires** your New Relic **region** to be set. You can set your region using the environment variable `NEW_RELIC_REGION` or by setting the `region` argument in your `provider` block. Valid values are `US` and `EU`.

    Using the environment variable:

    ```bash
    export NEW_RELIC_REGION="US" # or "EU"
    ```

    Using the `provider` block schema attribute:

    ```hcl
    provider "newrelic" {
      region = "US"
      # ... rest of config
    }
    ```

5. The v2 provider configuration **requires** a New Relic **account ID**. You can set your account ID using the environment variable `NEW_RELIC_ACCOUNT_ID` or by setting the `account_id` argument in your `provider` block.

    Using the environment variable:

    ```bash
    export NEW_RELIC_ACCOUNT_ID=<Your Account ID>
    ```

    Using the `provider` block:

    ```hcl
    provider "newrelic" {
      account_id = <Your Account ID>
      # ... rest of config
    }
    ```

#### API Key Environment Variables v1 to v2 Conversion Table

| v1                          | v2                        |
| --------------------------- | ------------------------- |
| `NEWRELIC_PERSONAL_API_KEY` | `NEW_RELIC_API_KEY`       |
| `NEWRELIC_API_KEY`          | Discontinued              |


### Provider Block Schema Updates

1. Replace any existing `api_key` configuration setting with the value of the existing `personal_api_key` configuration setting. The Personal API Key is now considered the default and standard API key for the provider.

    ```diff
    provider "newrelic" {
    -   api_key = "NRAA-***"
    -   personal_api_key="NRAK-***"
    +   api_key = "NRAK-***"

    }
    ```

    -> <small>**Note:** Take note of where the `NRAK-***` and `NRAA-***` prefixes switch. This is important. Most Personal API Keys have the `NRAK-` prefix.</small>

2. Add `account_id` to your `provider` block and set it to your New Relic account ID. Note that you can also use the environment variable `NEW_RELIC_ACCOUNT_ID`.

3. The `insights_account_id` configuration setting has been removed. The `account_id` configuration setting is now used instead.

[nr-personal-api-key-url]: https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#personal-api-key
