package main

import (
	"context"

	"github.com/machinebox/graphql"
)

func RunGraphQLRequest(client *graphql.Client, query string, variables map[string]interface{}, apiKey string, response interface{}) error {
	request := graphql.NewRequest(query)
	for key, value := range variables {
		request.Var(key, value)
	}
	request.Header.Set("Api-Key", apiKey)

	return client.Run(context.Background(), request, response)
}
