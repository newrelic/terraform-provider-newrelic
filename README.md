Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x

Using the provider
----------------------
To use a released provider in your Terraform environment, run `terraform init` and Terraform will automatically install the provider.

If you're developing and building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing the provider your plugins directory, run `terraform init` to initialize it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your
machine (version 1.13+ is *required*). You'll also need to correctly setup a
[GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

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
In order to test the provider, run `make test`. This will run the unit test suite.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests *create real resources*, and often cost money to run. The environment variables `NEWRELIC_API_KEY` and `NEWRELIC_LICENSE_KEY` must also be set with your associated keys for acceptance tests to work properly.

```sh
$ make testacc
```
