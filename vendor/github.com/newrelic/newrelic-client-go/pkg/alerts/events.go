package alerts

import (
	"time"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// AlertEvent response struct
type AlertEvent struct {
	ID            int                      `json:"id"`
	EventType     string                   `json:"event_type"`
	Product       string                   `json:"product"`
	EntityType    string                   `json:"entity_type"`
	EntityGroupID int                      `json:"entity_group_id"`
	EntityID      int                      `json:"entity_id"`
	Priority      string                   `json:"priority"`
	Description   string                   `json:"description"`
	Timestamp     *serialization.EpochTime `json:"timestamp"`
	IncidentID    int                      `json:"incident_id"`
}

// ListAlertEventsParams represents a set of filters to be used
// when querying New Relic alert events
type ListAlertEventsParams struct {
	Title         string     `url:"filter[title],omitempty"`
	Category      string     `url:"filter[category],omitempty"`
	CreatedAfter  *time.Time `url:"filter[created_after],omitempty"`
	CreatedBefore *time.Time `url:"filter[created_before],omitempty"`
	UpdatedAfter  *time.Time `url:"filter[updated_after],omitempty"`
	UpdatedBefore *time.Time `url:"filter[updated_before],omitempty"`
	Sort          string     `url:"sort,omitempty"`
	Page          int        `url:"page,omitempty"`
	PerPage       int        `url:"per_page,omitempty"`
}

// ListAlertEvents is used to retrieve New Relic alert events
func (a *Alerts) ListAlertEvents(params *ListAlertEventsParams) ([]*AlertEvent, error) {
	alertEvents := []*AlertEvent{}
	nextURL := a.config.Region().RestURL("alerts_events.json")

	for nextURL != "" {
		response := alertEventsResponse{}
		resp, err := a.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		alertEvents = append(alertEvents, response.AlertEvents...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertEvents, nil
}

type alertEventsResponse struct {
	AlertEvents []*AlertEvent `json:"alert_events,omitempty"`
}
