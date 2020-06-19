package main

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/region"
	nrProvider "github.com/terraform-providers/terraform-provider-newrelic/newrelic"
)

func main() {
	nrProvider.GetProvierUserAgentString(meta.SDKVersion)

	client := newClient()
	entities, err := client.Entities.SearchEntities(entities.SearchEntitiesParams{
		Domain: entities.EntityDomains.APM,
		Name:   "tf_test",
		// Reporting: &isReporting,
	})

	if err != nil {
		log.Fatalf("ERROR: %+v", err)
		return
	}

	if len(entities) > 0 {
		ch := make(chan int, len(entities))
		responses := []int{}
		deleted := []int{}

		for _, e := range entities {
			if e.ApplicationID != nil {
				// async
				go func(id int) {
					if _, err := client.APM.DeleteApplication(*e.ApplicationID); err != nil {
						log.Printf("[WARN] Error deleting application %v. Continuing to next entity...", err)
					} else {
						deleted = append(deleted, *e.ApplicationID)
					}

					ch <- *e.ApplicationID
				}(*e.ApplicationID)
			}
		}

		for {
			// do some stuff
			r, ok := <-ch
			if !ok {
				break
			}

			responses = append(responses, r)
			if len(responses) == len(entities) {
				log.Printf("deleted %d applications", len(deleted))
				close(ch)
			}
		}
	}
}

func newClient() *newrelic.NewRelic {
	client, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")),
		newrelic.ConfigAdminAPIKey(os.Getenv("NEW_RELIC_ADMIN_API_KEY")),
		newrelic.ConfigUserAgent(nrProvider.GetProvierUserAgentString(meta.SDKVersionString())),
		newrelic.ConfigServiceName("terraform-provider-newrelic"),
		newrelic.ConfigRegion(region.Name("US")),
	)

	if err != nil {
		log.Fatalf("failed to initialize client with error: %+v", err)
		return nil
	}

	return client
}
