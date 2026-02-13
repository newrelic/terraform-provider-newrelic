module github.com/newrelic/terraform-provider-newrelic/v3

go 1.25.0

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.26.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/newrelic/go-agent/v3 v3.30.0
	github.com/newrelic/go-insights v1.0.3
	github.com/newrelic/newrelic-client-go/v2 v2.77.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-git/go-git/v5 v5.12.0 // indirect
	github.com/go-test/deep v1.1.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.8 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hc-install v0.5.0 // indirect
	github.com/hashicorp/hcl/v2 v2.16.2 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.18.1 // indirect
	github.com/hashicorp/terraform-json v0.16.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.14.3 // indirect
	github.com/hashicorp/terraform-plugin-log v0.8.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.1.0 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jhump/protoreflect v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/robertkrimen/otto v0.5.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/zclconf/go-cty v1.13.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251222181119-0a764e51fe1b // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-viper/mapstructure/v2 v2.2.1 => github.com/go-viper/mapstructure/v2 v2.2.1

// Force remaining vulnerable dependencies to use fixed versions
replace (
	github.com/go-git/go-git/v5 => github.com/go-git/go-git/v5 v5.16.5
	github.com/sigstore/cosign/v2 => github.com/sigstore/cosign/v2 v2.6.2
	github.com/sigstore/rekor => github.com/sigstore/rekor v1.5.0
	github.com/sigstore/sigstore => github.com/sigstore/sigstore v1.10.4
)
