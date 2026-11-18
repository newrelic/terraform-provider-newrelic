//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenFleetEntity_Product_PCGFleet(t *testing.T) {
	t.Parallel()

	rd := resourceNewRelicFleet().TestResourceData()

	fleet := &fleetcontrol.EntityManagementFleetEntity{
		Name:              "my-pcg-fleet",
		ManagedEntityType: "KUBERNETESCLUSTER",
		Product:           []string{"PCG"},
	}

	err := flattenFleetEntity(fleet, rd, "org-123")
	require.NoError(t, err)

	product := rd.Get("product").([]interface{})
	assert.Len(t, product, 1)
	assert.Equal(t, "PCG", product[0].(string))
}

func TestFlattenFleetEntity_Product_NonPCGFleet(t *testing.T) {
	t.Parallel()

	rd := resourceNewRelicFleet().TestResourceData()

	fleet := &fleetcontrol.EntityManagementFleetEntity{
		Name:              "my-regular-fleet",
		ManagedEntityType: "KUBERNETESCLUSTER",
		Product:           nil,
	}

	err := flattenFleetEntity(fleet, rd, "org-123")
	require.NoError(t, err)

	product := rd.Get("product").([]interface{})
	assert.Empty(t, product, "non-PCG fleet should have empty product list")
}
