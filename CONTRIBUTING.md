# Contributing

At its core, this project is a Terraform provider, and so an understanding of
what Terraform is and how it works would help anyone looking to contribute to
this provider.

For those new to Terraform, the following might be good starting places.

-   [How Terraform Works][how_terraform_works]
-   [Plugin Types - Providers][how_terraform_works]

For those who are already familiar with Terraform, there are still plenty of
good resources on the [Extending Terraform][extending_terraform] page that are
worth looking at, but those might be a topic by topic basis, depending on the
need.  Also worth reading are the [Best Practices][best_practices].

## DTK Who?

The Developer Toolkit Team is a small team dedicated to open source, and
integrating New Relic APIs into the [Go Client][client_go], which we then
leverage in other projects, including this one.  You can read more about the
team [here][dtk].

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

### main.tf

When working to implement a feature or test a change, it's helpful to have an
HCL file that is isolated from any real production/staging environment.  When
using environment variables to [configure the provider][provider-config-guide],
the following is all that is needed to get a working provider for most
situations.

First export the necessary environment variables.

```bash
# Add this to your .bash_profile or .bashrc
export NEW_RELIC_API_KEY="<your New Relic Personal API key>"
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

### Changes to the provider

The simplest case is when the API calls you want to make are already in the [Go
Client][client_go] and the only changes you need to make are in the provider
code.  In that case, simply building the provider binary with the correct name
and running the plan will get you there.

```shell
export TARGET=.terraform/plugins/registry.terraform.io/newrelic/newrelic/99.0.0/darwin_amd64/
mkdir -p $TARGET
go build -o $TARGET/terraform-provider-newrelic && terraform init && TF_LOG=INFO terraform plan
```

The following assumes working on MacOS.  Essentially, we're making a "99.0.0"
fake version out of whatever we just built.  Then we can configure Terraform to
make use of the version we just built.

```hcl
terraform {
  required_providers {
    newrelic = {
      version = "~> 99.0.0"
      source = "newrelic/newrelic"
    }
  }
  required_version = ">= 0.13"
}
```

Now Terraform will expect the binary to exist in the `$TARGET` directory above.
See the [provider installation][provider_installation] and [provider
requirements][provider_requirements] docs for more information.

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
