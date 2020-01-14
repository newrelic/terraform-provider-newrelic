# newrelic-client-go

[![CircleCI](https://circleci.com/gh/newrelic/newrelic-client-go.svg?style=svg)](https://circleci.com/gh/newrelic/newrelic-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/newrelic/newrelic-client-go?style=flat-square)](https://goreportcard.com/report/github.com/newrelic/newrelic-client-go)
[![GoDoc](https://godoc.org/github.com/newrelic/newrelic-client-go?status.svg)](https://godoc.org/github.com/newrelic/newrelic-client-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/newrelic/newrelic-client-go/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/newrelic/newrelic-client-go?style=flat-square)](https://github.com/newrelic/newrelic-client-go/releases/latest)

The New Relic Client provides the building blocks for tools in the [Developer Toolkit](https://newrelic.github.io/developer-toolkit/), enabling quick access to the suite of New Relic APIs. As a library, it can also be leveraged within your own custom applications.

## Example

```go
import (
	"fmt"
	"os"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func main() {
	cfg := config.Config{
		APIKey: os.Getenv("NEWRELIC_API_KEY"),
	}

	nr := newrelic.New(cfg)

	params := apm.ListApplicationsParams{
		Name: "RPM",
	}

	apps, err := nr.APM.ListApplications(&params)

	if err != nil {
		fmt.Print(err)
	}
}

```


## Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. 

* [Roadmap](https://newrelic.github.io/developer-toolkit/roadmap/) - As part of the Developer Toolkit, the roadmap for this project follows the same RFC process
* [Issues or Enhancement Requests](https://github.com/newrelic/newrelic-client-go/issues) - Issues and enhancement requests can be submitted in the Issues tab of this repository. Please search for and review the existing open issues before submitting a new issue.
* [Contributors Guide](CONTRIBUTING.md) - Contributions are welcome (and if you submit a Enhancement Request, expect to be invited to contribute it yourself :grin:).
* [Community discussion board](https://discuss.newrelic.com/c/build-on-new-relic/developer-toolkit) - Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub.

Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. If you'd like to execute our corporate CLA, or if you have any questions, please drop us an email at opensource@newrelic.com.


## Development

### Requirements

* Go 1.13.0+
* GNU Make
* git


### Building

This package does not generate any direct usable assets (it's a library).  You can still run the build scripts to validate you code, and generate coverage information.

```
# Default target is 'build'
$ make

# Explicitly run build
$ make build

# Locally test the CI build scripts
# make build-ci
```


### Testing

Before contributing, all linting and tests must pass.  Tests can be run directly via:

```
# Tests and Linting
$ make test

# Only unit tests
$ make test-unit

# Only integration tests
$ make test-integration
```

### Commit Messages

Using the following format for commit messages allows for auto-generation of
the [CHANGELOG](CHANGELOG.md):

#### Format:

`<type>(<scope>): <subject>`

| Type | Description | Change log? |
|------| ----------- | :---------: |
| `chore` | Maintenance type work | No |
| `docs` | Documentation Updates | Yes |
| `feat` | New Features | Yes |
| `fix`  | Bug Fixes | Yes |
| `refactor` | Code Refactoring | No |

#### Scope

This refers to what part of the code is the focus of the work.  For example:

**General:**

* `build` - Work related to the build system (linting, makefiles, CI/CD, etc)
* `release` - Work related to cutting a new release

**Package Specific:**

* `newrelic` - Work related to the New Relic package
* `http` - Work related to the `internal/http` package
* `alerts` - Work related to the `pkg/alerts` package



### Documentation

**Note:** This requires the repo to be in your GOPATH [(godoc issue)](https://github.com/golang/go/issues/26827)

```
$ make docs
```


## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to the project here on GitHub.

_Please do not report issues with this software to New Relic Global Technical Support._


## Open Source License

This project is distributed under the [Apache 2 license](LICENSE).
