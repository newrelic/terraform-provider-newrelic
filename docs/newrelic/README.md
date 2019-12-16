# newrelic
--
    import "."


## Usage

```go
var (
	// Version is your app version (updated by Makefile, don't forget to TAG YOUR RELEASE)
	Version = "undefined"
)
```

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
