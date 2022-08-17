# Contributing

Contributions are always welcome. Before contributing please read the
[code of conduct](CODE_OF_CONDUCT.md) and [search the issue tracker](../../issues); your issue may have already been discussed or fixed in `main`. To contribute,
[fork](https://help.github.com/articles/fork-a-repo/) this repository, commit your changes, and [send a Pull Request](https://help.github.com/articles/using-pull-requests/).

Note that our [code of conduct](CODE_OF_CONDUCT.md) applies to all platforms and venues related to this project; please follow it in all your interactions with the project and its participants.

At its core, this project is a Terraform provider, and so an understanding of
what Terraform is and how it works would help anyone looking to contribute to
this provider.

For those new to Terraform, the following might be good starting places.

- [How Terraform Works][how_terraform_works]
- [Plugin Types - Providers][terraform_providers]

For those who are already familiar with Terraform, there are still plenty of
good resources on the [Extending Terraform][extending_terraform] page that are
worth looking at, but those might be a topic by topic basis, depending on the
need.  Also worth reading are the [Best Practices][best_practices].

## Future facing APIs

New Relic has several APIs, and it's worth understanding at a high level what
this means for the provider.  Also note, that each of the APIs for a given
product are the results of many teams, each owning their implementation,
documentation, and feature set, etc.

However, at a high level there is a concerted effort within New Relic to try
and move away from the REST APIs, and towards a newer GraphQL based API that we
call NerdGraph.  This migration is likely to take years.

What is mostly clear at this point, is that each team will continue to own
their migration, implementation, and feature set in terms of what they choose
to put into the GraphQL API.

### REST v2 APIs

Most of the documentation for the REST APIs are available at the [API
Explorer][api_explorer].  There you can see the call structures and responses,
and what methods are available.

### GraphQL API

To play and experiment with GraphQL, you can use the [graphiql][graphiql]
interface to perform queries and mutations.  Note that these are using the API,
so it's possible to make changes on your account.  Caution is advised when
performing mutations.

New fields, and new endpoints are added regularly in GraphQL, and so the
API surface area continues grow and scale.

## API Clients

There are a couple of different API clients that are used in this project, but
the bulk of the work is being done in the [New Relic Go Client][client_go].

There is also a Go Agent that is in use for testing, so that we have a mock
application to use during integration tests, and there is also an Insights
client that is likely on its way out of this project, so we won't go into
detail here.

## Development Process

### Building

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

### testing/newrelic.tf

When working to implement a feature or test a change, it's helpful to have an
HCL file that is isolated from any real production/staging environment.  When
using environment variables to [configure the provider][provider-config-guide],
the following is all that is needed to get a working provider for most
situations.

First export the necessary environment variables.

```bash
# Add this to your .bash_profile or .bashrc
export NEW_RELIC_API_KEY="<your New Relic User API key>"
export NEW_RELIC_ADMIN_API_KEY="<your New Relic Admin API key>"
export NEW_RELIC_REGION="US"
```

Then configure the provider.

```hcl
provider "newrelic" {}
```

Then you can begin to include a resource or data source that you want to test,
and experiment with changing attributes while running plan/apply to see the
results of how Terraform will behave.

### Testing

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

In order to run a single test, run the following command and replace `{testName}` with function name of your test.

```sh
TF_ACC=1 NR_ACC_TESTING=1 gotestsum -f testname -- -v --tags=integration -timeout 10m ./newrelic --run {testName}
```

### Changes to the provider

The simplest case is when the API calls you want to make are already in the [Go
Client][client_go] and the only changes you need to make are in the provider
code.  In that case, simply building the provider binary with the correct name
and running the plan will get you there. The following assumes working on MacOS.

To compile a new version of the compiler, run the following command in the root
directory of the repository. You will need to run this command every time you
make a code change.

```shell
make compile
```

To test your locally compiled plugin you can add your hcl files in the `testing`
directory. The `testing/dev.tfrc` file contains the local development configuration.
Before running any Terraform commands don't forgot to change the authentication
credentials in `testing/newrelic.tf` or use environment variables as mentioned above.
Additionally run the following command, or add it to your shell profile: `export TF_CLI_CONFIG_FILE=dev.tfrc`

You can now run `terraform plan` and
`terraform apply` in the `testing` directory to test your local version of the provider.

### Developing with our [Go Client][client_go]

When changes are required to the [Go Client][client_go] project, it might be
useful to test those changes in the client while developing the provider.  In
this case, you can tell Go to use your local copy of the [Go Client][client_go]
when building.

```shell
go mod edit -replace github.com/newrelic/newrelic-client-go=/Users/zleslie/go/src/github.com/newrelic/newrelic-client-go
```

This modifies the `go.mod` file to reference a local path on disk rather than
fetching the code from the remote `github.com/newrelic/newrelic-client-go`.
With this in place, you can then see how the two projects will behave when
there are changes in each that are needed for a particular feature.

Once complete, you can PR the change in the client repository, and then `git
checkout go.mod` in the provider to go back to a released version of client.
You'll want to make sure the version number in the provider lines up with the
client version you want to be using as well.

### Commit messages

To keep a style and allow us to automate the generating of a change log, we
require that that commit messages adhere to a standard.

TL;DR The commit message must match this regular expression.


  (chore|docs|feat|fix|refactor|tests?)(\([^\)]+\))?: .*


For more information on commit messages, we mostly follow [this standard][conventional_commits].

[api_explorer]: https://rpm.newrelic.com/api/explore/

[client_go]: https://github.com/newrelic/newrelic-client-go/

[dtk]: https://github.com/newrelic/developer-toolkit

[extending_terraform]: https://www.terraform.io/docs/extend/index.html

[graphiql]: https://api.newrelic.com/graphiql

[how_terraform_works]: https://www.terraform.io/docs/extend/how-terraform-works.html

[provider-config-guide]: https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/provider_configuration#configuration-via-environment-variables

[provider_design_principles]: https://www.terraform.io/docs/extend/hashicorp-provider-design-principles.html

[terraform_providers]: https://www.terraform.io/docs/extend/plugin-types.html#providers

[conventional_commits]: https://www.conventionalcommits.org/en/v1.0.0/

[best_practices]: https://www.terraform.io/docs/extend/best-practices/index.html

[provider_installation]: https://www.terraform.io/docs/commands/cli-config.html#provider-installation

[provider_requirements]: https://www.terraform.io/docs/configuration/terraform.html

## Feature Requests

Feature requests should be submitted in the [Issue tracker](../../issues), with a description of the expected behavior & use case, where they’ll remain closed until sufficient interest, [e.g. :+1: reactions](https://help.github.com/articles/about-discussions-in-issues-and-pull-requests/), has been [shown by the community](issues?q=label%3A%22votes+needed%22+sort%3Areactions-%2B1-desc).
Before submitting an Issue, please search for similar ones in the
[closed issues](issues?q=is%3Aissue+is%3Aclosed+label%3Aenhancement).

## Contributor License Agreement

Keep in mind that when you submit your Pull Request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.

For more information about CLAs, please check out Alex Russell’s excellent post,
[“Why Do I Need to Sign This?”](https://infrequently.org/2008/06/why-do-i-need-to-sign-this/).

## Slack

For contributors and maintainers of open source projects hosted by New Relic, we host a public Slack with a channel dedicated to this project. If you are contributing to this project, you're welcome to request access to that  community space.
