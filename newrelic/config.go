package newrelic

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

// Config contains New Relic provider settings
type Config struct {
	APIKey             string
	APIURL             string
	userAgent          string
	CACertFile         string
	InsecureSkipVerify bool
}

// Client returns a new client for accessing New Relic
func (c *Config) Client() (*newrelic.Client, error) {
	tlsCfg := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			log.Printf("Error reading CA Cert: %s", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsCfg.RootCAs = caCertPool
	}

	if c.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true
	}

	nrConfig := newrelic.Config{
		APIKey:    c.APIKey,
		BaseURL:   c.APIURL,
		Debug:     logging.IsDebugOrHigher(),
		TLSConfig: tlsCfg,
		UserAgent: c.userAgent,
	}

	client := newrelic.New(nrConfig)

	log.Printf("[INFO] New Relic client configured")

	return &client, nil
}

// ClientInfra returns a new client for accessing New Relic
func (c *Config) ClientInfra() (*newrelic.InfraClient, error) {
	nrConfig := newrelic.Config{
		APIKey:    c.APIKey,
		BaseURL:   c.APIURL,
		Debug:     logging.IsDebugOrHigher(),
		UserAgent: c.userAgent,
	}

	client := newrelic.NewInfraClient(nrConfig)

	log.Printf("[INFO] New Relic Infra client configured")

	return &client, nil
}

// ClientSynthetics returns a new client for accessing New Relic Synthetics
func (c *Config) ClientSynthetics() (*synthetics.Client, error) {
	conf := func(s *synthetics.Client) {
		s.APIKey = c.APIKey
	}

	client, _ := synthetics.NewClient(conf)

	log.Printf("[INFO] New Relic Synthetics client configured")

	return client, nil
}

// ProviderConfig for the custom provider
type ProviderConfig struct {
	Client      *newrelic.Client
	InfraClient *newrelic.InfraClient
	Synthetics  *synthetics.Client
}
