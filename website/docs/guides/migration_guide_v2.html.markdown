---
layout: "newrelic"
page_title: "New Relic Terraform Provider v2.x Migration Guide"
sidebar_current: "docs-newrelic-provider-v2-migration-guide"
description: |-
  Use this guide to update the New Relic Terraform Provider from v1.x to v2.x
---

## Update to v2.x of the New Relic Terraform Provider

Version 2.0 of the provider introduces some changes to the provider's configuration. Users wanting to upgrade from v1.x to v2.x will need to make a few adjustments to their configuration prior to upgrading.

### Environment Variable Updates

1. **IMPORTANT:** All environment variables in use by the provider have been renamed with a new naming convention. The `NEWRELIC_*` prefix has been changed to `NEW_RELIC_*`. This will be the environment variable naming convention moving forward.

2. The environment variable `NEWRELIC_PERSONAL_API_KEY` has been replaced with `NEW_RELIC_API_KEY`.

    ```diff
    - NEWRELIC_PERSONAL_API_KEY
    + NEW_RELIC_API_KEY
    ```

3. The environment variable `NEWRELIC_API_KEY` has been replaced with `NEW_RELIC_ADMIN_API_KEY`.

    ```diff
    - NEWRELIC_API_KEY
    + NEW_RELIC_ADMIN_API_KEY
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



### Provider Block Schema Updates

1. Move any existing `api_key` configuration setting to `admin_api_key`. Then move any existing `personal_api_key` configuration setting to `api_key`. Your Personal API Key is now considered the default and standard API key for the provider.

    ```diff
    provider "newrelic" {
    -   api_key = "NRAA-***"
    -   personal_api_key="NRAK-***"
    +   api_key = "NRAK-***"
    +   admin_api_key = "NRAA-***"

    }
    ```

    <small>**Note:** Take note of where the `NRAK-***` and `NRAA-***` values switch. This is important. Your Personal API Key has the `NRAK-` prefix.</small>

2. Add `account_id` to your `provider` block.

3. The `insights_account_id` configuration setting has been removed. The `account_id` configuration setting is now used instead.



<br><br><br><br><br><br><br>
-> <small>**Note:** Make sure to double check your different API keys are set the correct environment variables. <br> `NEW_RELIC_API_KEY` should be set to your Personal API Key (formatted `NRAK-***`). <br>`NEW_RELIC_ADMIN_API_KEY` should be set to your Admin API Key (formatted `NRAA-***`).</small>


**Example v1 naming convention comparison:**

| v1                 | v2                  |
| ------------------ | ------------------- |
| `NEWRELIC_API_KEY` | `NEW_RELIC_API_KEY` |
