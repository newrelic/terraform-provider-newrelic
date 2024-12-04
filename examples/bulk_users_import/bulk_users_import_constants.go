package main

const GET_GROUPS_AND_USERS_QUERY = `
query(
	$groupId: [ID!]
){
	actor {
		organization {
			userManagement {
				authenticationDomains {
				nextCursor
					authenticationDomains {
						name
						id
						groups(id: $groupId) {
							groups {
								id
								displayName
								users {
									users {
										email
										id
										name
										timeZone
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
`

const GET_USER_DETAILS_QUERY = `
query(
	$userId: [ID!]
){
	actor {
		organization {
			userManagement {
				authenticationDomains {
					nextCursor
					authenticationDomains {
						name
						id
						users(id: $userId) {
							users {
								email
								id
								name
								type {
									displayName
								}
							}
						}
					}
				}
			}
		}
	}
}`

const NERDGRAPH_API_ENDPOINT = "https://api.newrelic.com/graphql"
const NERDGRAPH_API_ENDPOINT_EU = "https://api.newrelic.com/graphql"
