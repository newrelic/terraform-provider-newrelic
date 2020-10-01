---
name: Bug report
about: Report a bug
title: ''
labels: ''
assignees: ''

---

Hi there,

Thank you for opening an issue. In order to better assist you with your issue, we kindly ask to follow the template format and instructions. Please note that we try to keep the Terraform issue tracker reserved for **bug reports** and **feature requests** only. General usage questions submitted as issues will be closed and redirected to New Relic's Explorers Hub https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit.

## Please include the following with your bug report

> :warning: **Important:** Failure to include the following, such as omitting the Terraform configuration in question, may delay resolving the issue.

- [ ] Your New Relic `provider` [configuration](#terraform-configuration-files) (sensitive details redacted)
- [ ] A list of [affected resources](#affected-resources) and/or data sources
- [ ] The [configuration](#terraform-configuration-files) of the resources and/or data sources related to the bug report (i.e. from the list mentioned above)
- [ ] Description of the [current behavior](#actual-behavior) (the bug)
- [ ] Description of the [expected behavior](#expected-behavior)
- [ ] Any related [log output](#debug-output)


### Terraform Version
Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

### Affected Resource(s)
Please list the resources as a list, for example:
- `newrelic_alert_policy`
- `newrelic_alert_channel`

If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this.

### Terraform Configuration
> Please include your `provider` configuration (sensitive details redacted) as well as the configuration of the resources and/or data sources related to the bug report.
```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.
```

### Actual Behavior
What actually happened?

### Expected Behavior
What should have happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### Debug Output
Please provider a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

### Panic Output
If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`.

### Important Factoids
Are there anything atypical about your accounts that we should know? For example: Running in EC2 Classic? Custom version of OpenStack? Tight ACLs?

### References
Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- GH-1234
