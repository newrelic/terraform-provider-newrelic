package main

import (
	"fmt"
	"log"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

func main() {
	client := NewClient()

	entities, err := client.Entities.SearchEntities(entities.SearchEntitiesParams{
		Domain: entities.EntityDomains.APM,
		Name:   "tf_test",
		// Reporting: &isReporting,
	})

	if err != nil {
		log.Fatalf("ERROR: %+v", err)
		return
	}

	fmt.Print("\n****************************\n")

	if len(entities) > 0 {
		for _, e := range entities {
			if e.ApplicationID != nil {
				if _, err := client.APM.DeleteApplication(*e.ApplicationID); err != nil {
					log.Printf("ERROR: failed to delete application with ID %v due to error: %v . Continuing to next entity...", *e.ApplicationID, err)
					continue
				}

				log.Printf("successfully deleted application with ID: %+v\n", *e.ApplicationID)
			}
		}
	}
}

func NewClient() *newrelic.NewRelic {
	client, err := newrelic.New(
		// newrelic.ConfigPersonalAPIKey(os.Getenv("NEW_RELIC_API_KEY")),
		// newrelic.ConfigAdminAPIKey(os.Getenv("NEW_RELIC_ADMIN_API_KEY")),
		// newrelic.ConfigUserAgent(os.Getenv("")),
		// newrelic.ConfigServiceName(os.Getenv("")),
		newrelic.ConfigRegion(region.Name("US")),
	)

	if err != nil {
		log.Fatalf("failed to initialize client with error: %+v", err)
		return nil
	}

	return client
}
