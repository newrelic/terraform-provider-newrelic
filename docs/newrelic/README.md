# newrelic
--
    import "."


## Usage

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
}
```

Config contains all the configuration data for the API Client.

#### func (*Config) ToInternal

```go
func (c *Config) ToInternal() internal.Config
```
