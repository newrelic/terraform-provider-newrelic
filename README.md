# New Relic Terraform Provider

- Documentation: <https://registry.terraform.io/providers/newrelic/newrelic/latest/docs>
- Terraform Website: <https://www.terraform.io>
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.11+

## Using the provider

To use the latest version of the provider in your Terraform environment, run `terraform init` and Terraform will automatically install the provider.

If you wish to pin your environment to a specific release of the provider, you can do so with a `required_providers` statement in your Terraform manifest. The `terraform` [configuration block](https://www.terraform.io/docs/configuration/provider-requirements.html) varies slightly depending on which Terraform version you're using. See below for more examples of configuring the provider version for the different versions of Terraform.

For Terraform version 0.13.x

```hcl
terraform {
  required_version = "~> 0.13.0"
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 2.7.5"
    }
  }
}
```

For Terraform version 0.12.x

```hcl
terraform {
  required_providers {
    newrelic = {
      version = "~> 2.7.5"
    }
  }
}
```

For Terraform version 0.11.x

```hcl
provider "newrelic" {
  version = "~> 2.7.5"
}
```

If you're developing and building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing the provider your plugins directory, run `terraform init` to initialize it.

For more information on using the provider and the associated resources, please see the [provider documentation][provider_docs] page.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your
machine (version 1.13+ is _required_). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Please see our [CONTRIBUTING][contributing] guide for more detail on the APIs
in use by this provider.

#### Building

Clone repository to: `$GOPATH/src/github.com/newrelic/terraform-provider-newrelic`

```sh
$ mkdir -p $GOPATH/src/github.com/newrelic;
$ cd $GOPATH/src/github.com/newrelic
$ git clone git@github.com:newrelic/terraform-provider-newrelic.git
```

Enter the provider directory and build the provider. To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ cd $GOPATH/src/github.com/newrelic/terraform-provider-newrelic
$ make build
```

#### Testing

In order to test the provider, run `make test`. This will run the full test suite.

```sh
$ make test
```

In order to run the unit test suite only, run `make test-unit`.

```sh
$ make test-unit
```

In order to run the acceptance test suite only, run `make test-integration`.

```sh
$ make test-integration
```

_Note:_ Acceptance tests _create real resources_. The following environment
variables must bet set for acceptance tests to run:

```sh
NEW_RELIC_API_KEY
NEW_RELIC_ACCOUNT_ID
NEW_RELIC_INSIGHTS_INSERT_KEY
NEW_RELIC_LICENSE_KEY
NEW_RELIC_REGION
```

#### Go Version Support

We'll aim to support the latest supported release of Go, along with the
previous release. This doesn't mean that building with an older version of Go
will not work, but we don't intend to support a Go version in this project that
is not supported by the larger Go community. Please see the [Go
releases][go_releases] page for more details.

[provider_docs]: https://www.terraform.io/docs/providers/newrelic/index.html
[contributing]: https://github.com/newrelic/terraform-provider-newrelic/blob/main/CONTRIBUTING.md
[go_releases]: https://github.com/golang/go/wiki/Go-Release-Cycle
