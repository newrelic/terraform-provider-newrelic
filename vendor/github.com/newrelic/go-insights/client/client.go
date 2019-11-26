package client

import (
	"net/url"
)

// UseCustomURL allows overriding the default Insights Host / Scheme.
func (c *Client) UseCustomURL(customURL string) {
	newURL, _ := url.Parse(customURL)
	if len(newURL.Scheme) < 1 {
		c.URL.Scheme = "https"
	} else {
		c.URL.Scheme = newURL.Scheme
	}

	c.URL.Host = newURL.Host
	c.Logger.Debugf("Using custom URL: %s", c.URL)
}
