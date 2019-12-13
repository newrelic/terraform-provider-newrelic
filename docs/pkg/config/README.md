# config
--
    import "."


## Usage

```go
const (
	// Production represents New Relic's US-based production deployment.
	Production = iota

	// EU represents New Relic's EU-based production deployment.
	EU

	// Staging represents New Relic's US-based staging deployment.  This is for internal use only.
	Staging
)
```

#### type Config

```go
type Config struct {
	APIKey        string
	BaseURL       string
	ProxyURL      string
	Debug         bool
	TLSConfig     *tls.Config
	UserAgent     string
	HTTPTransport http.RoundTripper
	Region        Region
}
```

Config contains all the configuration data for the API Client.

#### type Region

```go
type Region int
```

Region specifies the New Relic environment to target.
