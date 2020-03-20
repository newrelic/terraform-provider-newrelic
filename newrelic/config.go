package newrelic

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	insights "github.com/newrelic/go-insights/client"
	nr "github.com/newrelic/newrelic-client-go/newrelic"
)

const serviceName = "terraform-provider-newrelic"

// Config contains New Relic provider settings
type Config struct {
	AdminAPIKey          string
	PersonalAPIKey       string
	APIURL               string
	CACertFile           string
	InfrastructureAPIURL string
	InsecureSkipVerify   bool
	InsightsAccountID    string
	InsightsInsertKey    string
	InsightsInsertURL    string
	InsightsQueryKey     string
	InsightsQueryURL     string
	NerdGraphAPIURL      string
	SyntheticsAPIURL     string
	userAgent            string
}

// Client returns a new client for accessing New Relic
func (c *Config) Client() (*nr.NewRelic, error) {
	options := []nr.ConfigOption{}

	options = append(options,
		nr.ConfigAdminAPIKey(c.AdminAPIKey),
		nr.ConfigPersonalAPIKey(c.PersonalAPIKey),
		nr.ConfigUserAgent(c.userAgent),
		nr.ConfigServiceName(serviceName),
	)

	tlsCfg := &tls.Config{}
	var t = http.DefaultTransport

	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			log.Printf("Error reading CA Cert: %s", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsCfg.RootCAs = caCertPool

		t = &http.Transport{TLSClientConfig: tlsCfg}
	} else if c.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true

		t = &http.Transport{TLSClientConfig: tlsCfg}
	}

	if logging.LogLevel() != "" {
		options = append(options, nr.ConfigLogLevel(logging.LogLevel()))
		t = logging.NewTransport("newrelic", t)
	}

	options = append(options, nr.ConfigHTTPTransport(t))

	if c.APIURL != "" {
		options = append(options, nr.ConfigBaseURL(c.APIURL))
	}

	if c.SyntheticsAPIURL != "" {
		options = append(options, nr.ConfigSyntheticsBaseURL(c.SyntheticsAPIURL))
	}

	if c.InfrastructureAPIURL != "" {
		options = append(options, nr.ConfigInfrastructureBaseURL(c.InfrastructureAPIURL))
	}

	if c.NerdGraphAPIURL != "" {
		options = append(options, nr.ConfigNerdGraphBaseURL(c.NerdGraphAPIURL))
	}

	client, err := nr.New(options...)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] New Relic client configured")

	return client, nil
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

// ProviderConfig for the custom provider
type ProviderConfig struct {
	NewClient            *nr.NewRelic
	InsightsInsertClient *insights.InsertClient
	InsightsQueryClient  *insights.QueryClient
}
