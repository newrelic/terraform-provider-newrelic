// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListLabelsResponseJSON = `{
		"labels": [
			{
				"key": "Project:example-project",
				"category": "Project",
				"name": "example-project",
				"application_health_status": {
					"green": [],
					"orange": [],
					"red": [],
					"gray": []
				},
				"server_health_status": {
					"green": [],
					"orange": [],
					"red": [],
					"gray": []
				},
				"origins": {
					"apm": [],
					"synthetics": [],
					"agents": [
						12345
					]
				},
				"links": {
					"applications": [
						12345
					],
					"servers": []
				}
			}
		]
	}`
)

func TestListLabels(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testListLabelsResponseJSON, http.StatusOK)

	expected := []*Label{
		{
			Key:      "Project:example-project",
			Category: "Project",
			Name:     "example-project",
			Links: LabelLinks{
				Applications: []int{12345},
				Servers:      []int{},
			},
		},
	}

	actual, err := apm.ListLabels()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetLabel(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testListLabelsResponseJSON, http.StatusOK)

	expected := &Label{
		Key:      "Project:example-project",
		Category: "Project",
		Name:     "example-project",
		Links: LabelLinks{
			Applications: []int{12345},
			Servers:      []int{},
		},
	}

	actual, err := apm.GetLabel("Project:example-project")

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
