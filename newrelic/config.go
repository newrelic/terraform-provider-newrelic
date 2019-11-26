package newrelic

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/url"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	insights "github.com/newrelic/go-insights/client"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

// Config contains New Relic provider settings
type Config struct {
	InsightsAccountID  string
	APIKey             string
	APIURL             string
	InsightsInsertKey  string
	InsightsInsertURL  string
	InsightsQueryKey   string
	InsightsQueryURL   string
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
		UserAgent: c.userAgent,
		TLSConfig: tlsCfg,
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

// ClientInsightsInsert returns a new Insights insert client
func (c *Config) ClientInsightsInsert() (*insights.InsertClient, error) {
	client := insights.NewInsertClient(c.InsightsInsertKey, c.InsightsAccountID)

	if c.InsightsInsertURL != "" {
		insightsURL, err := url.Parse(c.InsightsInsertURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing Insights URL: %q", err)
		}
		insightsURL.Path = fmt.Sprintf("%s/%s/events", insightsURL.Path, c.InsightsAccountID)
		client.URL = insightsURL
	}

	client.SetCompression(gzip.DefaultCompression)

	if len(c.InsightsInsertKey) > 1 {
		if err := client.Validate(); err != nil {
			return nil, err
		}
	}

	log.Printf("[INFO] New Relic Insights insert client configured")

	return client, nil
}

// ClientInsightsQuery returns a new Insights query client
func (c *Config) ClientInsightsQuery() (*insights.QueryClient, error) {
	client := insights.NewQueryClient(c.InsightsQueryKey, c.InsightsAccountID)

	if c.InsightsQueryURL != "" {
		insightsURL, err := url.Parse(c.InsightsQueryURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing Insights URL: %q", err)
		}
		insightsURL.Path = fmt.Sprintf("%s/%s/query", insightsURL.Path, c.InsightsAccountID)
		client.URL = insightsURL
	}

	if len(c.InsightsQueryKey) > 1 {
		if err := client.Validate(); err != nil {
			return nil, err
		}
	}

	log.Printf("[INFO] New Relic Insights query client configured")

	return client, nil
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
	Client               *newrelic.Client
	InsightsInsertClient *insights.InsertClient
	InsightsQueryClient  *insights.QueryClient
	InfraClient          *newrelic.InfraClient
	Synthetics           *synthetics.Client
}
