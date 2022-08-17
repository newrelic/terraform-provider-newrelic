[![Community Plus header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Plus.png)](https://opensource.newrelic.com/oss-category/#community-plus)

# New Relic Terraform Provider

New Relic offers tools that help you fix problems quickly, maintain complex systems, improve your code, and accelerate your digital transformation. With the New Relic Terraform provider you are able to automate the configuration of New Relic.

- Documentation: <https://registry.terraform.io/providers/newrelic/newrelic/latest/docs>
- Terraform Website: <https://www.terraform.io>
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0+

New Relic and the Terraform team will support Terraform versions up to 2 years after the latest release. We advice to always upgrade to the latest version of Terraform and the New Relic Terraform provider.

## Using the provider

To use the latest version of the provider in your Terraform environment, run `terraform init` and Terraform will automatically install the provider.

If you wish to pin your environment to a specific release of the provider, you can do so with a `required_providers` statement in your Terraform manifest. The `terraform` [configuration block](https://www.terraform.io/docs/configuration/provider-requirements.html) varies slightly depending on which Terraform version you're using. See below for more examples of configuring the provider version for the different versions of Terraform.

For Terraform version 1.x and above

```hcl
terraform {
  required_version = "~> 1.0"
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}
```

If you're developing and building the provider locally, follow the [instructions in our contribution guide](https://github.com/newrelic/terraform-provider-newrelic/blob/main/CONTRIBUTING.md#development-process).

For more information on using the provider and the associated resources, please see the [provider documentation][provider_docs] page.

## Support

Should you need assistance with New Relic products, you are in good hands with several support channels.

**Support Channels**

* [New Relic Documentation](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs): Comprehensive guidance for using our Terraform provider
* [New Relic Community](https://discuss.newrelic.com/tag/terraform): The best place to engage in troubleshooting questions
* [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
* [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level
* [New Relic Technical Support](https://support.newrelic.com/) 24/7/365 ticketed support. Read more about our [Technical Support Offerings](https://docs.newrelic.com/docs/licenses/license-information/general-usage-licenses/global-technical-support-offerings).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your
machine (version 1.18 is _required_). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Please see our [CONTRIBUTING][contributing] guide for more information about developing and testing the New Relic Terraform provider.

#### Go Version Support

We'll aim to support the latest supported release of Go, along with the
previous release. This doesn't mean that building with an older version of Go
will not work, but we don't intend to support a Go version in this project that
is not supported by the larger Go community. Please see the [Go releases][go_releases] page for more details.

[provider_docs]: https://www.terraform.io/docs/providers/newrelic/index.html
[contributing]: https://github.com/newrelic/terraform-provider-newrelic/blob/main/CONTRIBUTING.md
[go_releases]: https://github.com/golang/go/wiki/Go-Release-Cycle

