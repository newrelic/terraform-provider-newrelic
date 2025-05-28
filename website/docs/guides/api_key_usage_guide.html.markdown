---
layout: "newrelic"
page_title: "Securely Providing Your New Relic API Key to Terraform"
sidebar_current: "docs-newrelic-provider-api-key-usage-guide"
description: |-
  Use this guide to fully automate your New Relic setup and configurations through Terraform, with your API key securely managing the programmatic access.
---

# Important: Securely Providing Your New Relic API Key to Terraform

To prevent the exposure of your sensitive New Relic API Key in code repositories, state files, or logs, you must use an environment variable to provide this credential to Terraform.
* **Avoid Hardcoding:** Never hardcode your API key directly in `.tf` files. These files are often committed to version control (like Git), making your sensitive key publicly visible and vulnerable.
* **Prevent Exposure:** Standard Terraform plan/apply outputs can sometimes expose variable values. Using environment variables helps keep them out of these logs.
* **CI/CD Best Practice:** In automated environments (e.g., GitHub Actions, GitLab CI, Jenkins), API keys should always be managed as secrets and injected via environment variables.

---

## How to Securely Provide Your API Key (Mandatory Method)

Follow these steps to ensure your `newrelic_api_key ` is handled securely:

**Define the Variable in your Terraform Configuration:**
Ensure you have the following `variable` block in a `.tf` file within your root module (e.g., `variables.tf`):

```text
variable "newrelic_api_key" {
  type        = string
  sensitive   = true # Marks the variable as sensitive for masking in outputs
  default     = null # Ensures no interactive prompt if not set, making it CI/CD friendly
}
```

**Set the Environment Variable:**
You **must** set an environment variable named `TF_VAR_newrelic_api_key` (note the `TF_VAR_` prefix and the uppercase variable name matching your `variable "newrelic_api_key"` block). Terraform will automatically read this value.

* **Local Development (before running `terraform plan` or `apply`):**

    * Linux/macOS
   ```powershell
    export TF_VAR_newrelic_api_key="NRAK-YOUR-ACTUAL-NEW-RELIC-API-KEY" 
    ```
    * Windows Command Prompt
    ```powershell
    set TF_VAR_newrelic_api_key="NRAK-YOUR-ACTUAL-NEW-RELIC-API-KEY"
    ```
    * Windows PowerShell
    ```powershell
    $env:TF_VAR_newrelic_api_key ="NRAK-YOUR-ACTUAL-NEW-RELIC-API-KEY"
    ```
* **CI/CD Pipelines (Recommended Method):**
  Always use your CI/CD platform's **secret management features** to securely store your API key. Then, expose this secret as an environment variable named `TF_VAR_newrelic_api_key` to your Terraform job.

* **Example (GitHub Actions):**

    ```yaml
    # ... within your workflow job step ...
    - name: Terraform Plan
      env:
        TF_VAR_newrelic_api_key: ${{ secrets.NEW_RELIC_API_KEY }} # Assumes 'NEW_RELIC_API_KEY' is a GitHub Secret
      run: terraform plan
    ```
