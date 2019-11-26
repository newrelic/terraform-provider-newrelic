package newrelic

import (
	"testing"
)

func TestConfigClient_Basic(t *testing.T) {
	config := &Config{
		APIKey: "foo",
	}

	nr, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	if nr == nil || nr.RestyClient == nil {
		t.Fatal("failed to create newrelic client")
	}
	if nr.RestyClient.HostURL != "https://api.newrelic.com/v2" {
		t.Fatalf("failed to set default Client APIURL. expected %v, got %v", "https://api.newrelic.com/v2", nr.RestyClient.HostURL)
	}
}

func TestConfigClient_CACertFile(t *testing.T) {
	config := &Config{
		APIKey:     "foo",
		CACertFile: "notafile",
	}

	_, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigClient_InsecureSkipVerify(t *testing.T) {
	config := &Config{
		APIKey:             "foo",
		InsecureSkipVerify: true,
	}

	_, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigClientInfra_Basic(t *testing.T) {
	config := &Config{
		APIKey: "foo",
		APIURL: "https://infra-api.newrelic.com/v2",
	}

	nr, err := config.ClientInfra()
	if err != nil {
		t.Fatal(err)
	}
	if nr == nil || nr.RestyClient == nil {
		t.Fatal("failed to create newrelic client")
	}
	if nr.RestyClient.HostURL != config.APIURL {
		t.Fatalf("failed to set default ClientInfra APIURL. expected %v, got %v", config.APIURL, nr.RestyClient.HostURL)
	}
}

func TestConfigClientSynthetics_Basic(t *testing.T) {
	config := &Config{
		APIKey: "foo",
	}

	_, err := config.ClientSynthetics()
	if err != nil {
		t.Fatal(err)
	}
}
