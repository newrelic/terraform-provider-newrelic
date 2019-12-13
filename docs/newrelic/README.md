# newrelic
--
    import "."


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
func New(config config.Config) NewRelic
```
