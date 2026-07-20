// Package newrelic comprises the implementation of all resources in the New Relic Terraform Provider
package newrelic

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mitchellh/go-homedir"

	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
)

// Config contains New Relic provider settings
type Config struct {
	AdminAPIKey          string
	PersonalAPIKey       string
	Region               string
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
	serviceName          string
}

// Client returns a new client for accessing New Relic
func (c *Config) Client() (*nr.NewRelic, error) {
	options := []nr.ConfigOption{}

	options = append(options,
		nr.ConfigAdminAPIKey(c.AdminAPIKey),
		nr.ConfigPersonalAPIKey(c.PersonalAPIKey),
		nr.ConfigInsightsInsertKey(c.InsightsInsertKey),
		nr.ConfigUserAgent(c.userAgent),
		nr.ConfigServiceName(c.serviceName),
		nr.ConfigRegion(c.Region),
	)

	tlsCfg := &tls.Config{}
	var t = http.DefaultTransport

	if c.CACertFile != "" {
		caCert, _, err := read(c.CACertFile)
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

// ProviderConfig for the custom provider
type ProviderConfig struct {
	NewClient      *nr.NewRelic
	AccountID      int
	PersonalAPIKey string
	Region         string
	userAgent      string
}

func (p *ProviderConfig) GetUserAgent() string {
	return p.userAgent
}

// If the argument is a path, Read loads it and returns the contents,
// otherwise the argument is assumed to be the desired contents and is simply
// returned.
//
// The boolean second return value can be called `wasPath` - it indicates if a
// path was detected and a file loaded.
func read(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}
