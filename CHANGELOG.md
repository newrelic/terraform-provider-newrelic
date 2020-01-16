<a name="unreleased"></a>
## [Unreleased]


<a name="v0.5.0"></a>
## [v0.5.0] - 2020-01-16
### Documentation Updates
- **newrelic:** update API key configuration documentation


<a name="v0.4.0"></a>
## [v0.4.0] - 2020-01-15
### Bug Fixes
- retry HTTP requests on 429 status codes

### Features
- **entities:** add entities search and entity tagging


<a name="v0.3.0"></a>
## [v0.3.0] - 2020-01-13
### Bug Fixes
- make use of ErrorNotFound type for Get methods that are based on List methods
- add policy ID to alert condition

### Documentation Updates
- update example
- **build:** Update README for commit message format
- **changelog:** Add auto-generation of CHANGELOG from git comments via `make changelog`

### Features
- add top-level logging package for convenience
- add option for JSON logging and fail gracefully when log level cannot be parsed
- introduce logging
- update monitor scripts with return design pattern, update tests


<a name="v0.2.0"></a>
## [v0.2.0] - 2020-01-08
### Documentation Updates
- update readme example


<a name="v0.1.0"></a>
## v0.1.0 - 2020-01-07
### Bug Fixes
- rename variables to fix redeclared error
- update unit tests to use new method sigs
- fix monitor ID type and GetMonitor URL
- http client needs to handle other 'success' response status codes such as 201
- add godoc as a dep, and a warning about GOPATH and godoc
- fix paging bug for v2 API
- **lint:** formatting fixes for linter

### Documentation Updates
- add alerts package docs
- temporarily checking in broken import paths in generated markdown docs
- add inline documentation
- add badges to README
- fill in missing inline documentation
- document some methods

### Features
- add DeletePluginCondition
- add CreatePluginCondition
- add UpdatePluginCondition
- add GetPluginCondition
- add ListPluginsConditions
- encode monitor script text
- add ability to use 'detailed' query param in ListPlugins method
- add GetPlugin
- add ListPlugins
- publicly expose error types
- finish components endpoints
- add Components
- add internal utils package, move IntArrayToString() util to new home
- add integration tests for key transactions
- add query param filters for ListKeyTransactions
- add GetKeyTransaction
- add ListKeyTransactions
- add DeleteLabel
- add CreateLabel
- add ListLabels, add GetLabel
- add DeleteDeployment
- add CreateDeployment
- add ListDeployments
- centralize apm test helpers
- add DeleteNrqlAlertCondition
- add UpdateNrqlAlertCondition
- add CreateNrqlAlertCondition
- add GetNrqlAlertCondition
- add ListNrqlAlertConditions
- add UpdateAlertPolicy
- add DeleteAlertCondition
- add CreateAlertCondition
- add GetAlertCondition
- add ListAlertConditions
- get infra condition integration tests passing
- add InfrastructureConditions
- add MonitorScripts
- add MonitorScript
- add DeleteAlertPolicyChannel, update unit tests, add integration test (might need to remove this)
- add alert policy channels
- add synthetics alert conditions
- add synthetics alert conditions
- add GetAlertChannel method
- add CreateAlertChannel, ListAlertChannels, DeleteAlertChannel
- add DeleteMonitor
- add UpdateMonitor
- add CreateMonitor
- add dashboards
- add DeleteAlertPolicy method
- add UpdateAlertPolicy method
- add CreateAlertPolicy method
- add GetAlertPolicy method
- add ListAlertPolicies method
- alerts package
- create remaining CRUD methods for application resource
- add new dependency-free client implementation
- add version.go per auto-versioning docs
- add ListAlertConditions for infrastructure
- add infra namespace
- add catchall newrelic package
- add New Relic environment enum
- maximize page size for ListMonitors
- add ListMonitors method for Synthetics monitors
- add application filtering for ListApplications
- get TestListApplications passing


[Unreleased]: https://github.com/newrelic/newrelic-client-go/compare/v0.5.0...HEAD
[v0.5.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/newrelic/newrelic-client-go/compare/v0.1.0...v0.2.0
