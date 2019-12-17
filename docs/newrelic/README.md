# newrelic
--
    import "github.com/newrelic/newrelic-client-go/newrelic"


## Usage

#### type NewRelic

```go
type NewRelic struct {
	APM            apm.APM
	Synthetics     synthetics.Synthetics
	Infrastructure infrastructure.Infrastructure
}
```


#### func  New

```go
func New(config config.ReplacementConfig) NewRelic
```
