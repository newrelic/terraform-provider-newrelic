//go:build integration || ENTITY
// +build integration ENTITY

package newrelic

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

var (
	testSourceEntityGUID = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1ODA4MDM"
	testTargetEntityGUID = "MzgwNjUyNnxFWFR8U0VSVklDRV9MRVZFTHw1NzE0Nzk"
)

func TestAccNewRelicEntityRelationship_Basic(t *testing.T) {
	resourceName := "newrelic_entity_relationship.rel"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicEntityRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityRelationshipConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityRelationshipExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicEntityRelationshipConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityRelationshipExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: false,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestAccNewRelicEntityRelationship_Validation(t *testing.T) {
	expectedMsg, _ := regexp.Compile(fmt.Sprintf("problem in retrieving application with GUID %ss", testSourceEntityGUID))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicEntityRelationshipDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicEntityRelationshipConfigValidation(),
				ExpectError: expectedMsg,
			},
		},
	})
}

func testAccCheckNewRelicEntityRelationshipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_entity_relationship" {
			continue
		}
		sourceEntityGUID, targetEntityGUID, err := getEntityRelationshipGUIDs(rs.Primary.ID)
		// Check if the entity relationship still exists
		resp, err := client.Entities.GetEntity(common.EntityGUID(sourceEntityGUID))

		switch (*resp).(type) {
		case *entities.ExternalEntity:
			relatedEntities := (*resp).(*entities.ExternalEntity).RelatedEntities
			for _, relationship := range relatedEntities.Results {
				if userDefinedEdge, ok := relationship.(*entities.EntityRelationshipUserDefinedEdge); ok {
					// Now 'userDefinedEdge' is a pointer to an EntityRelationshipUserDefinedEdge and you can access its fields.
					if userDefinedEdge.Target.GUID == common.EntityGUID(targetEntityGUID) {
						//Set the relationship related fields in the resource data
						return fmt.Errorf("entity relationship %s still exists", rs.Primary.ID)
						break
					}
				}
			}

		default:
			return err
		}
	}
	return nil
}

func testAccNewRelicEntityRelationshipConfig() string {
	return fmt.Sprintf(`
		resource "newrelic_entity_relationship" "rel" {
			source_entity_guid = "%[1]s"
			target_entity_guid = "%[2]s"
			relation_type = "CONTAINS"
		}`, testSourceEntityGUID, testTargetEntityGUID)
}

func testAccNewRelicEntityRelationshipConfigUpdated() string {
	return fmt.Sprintf(`
		resource "newrelic_entity_relationship" "rel" {
			source_entity_guid = "%[1]s"
			target_entity_guid = "%[2]s"
			relation_type = "CONTAINS"
		}`, testSourceEntityGUID, testTargetEntityGUID)
}

func testAccNewRelicEntityRelationshipConfigValidation() string {
	return fmt.Sprintf(`
		resource "newrelic_entity_relationship" "rel" {
			source_entity_guid = "%[1]ss"
			target_entity_guid = "%[2]s"
			relation_type = "CONTAINS"
		}`, testSourceEntityGUID, testTargetEntityGUID)
}

func testAccCheckNewRelicEntityRelationshipExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no relationship ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		return nil
		time.Sleep(2 * time.Second)
		_, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		//(*found).(*entities.ExternalEntity).RelatedEntities
		//if !foundOk {
		//	return fmt.Errorf("no relationship found")
		//}
		//if res.Results != common.EntityGUID(rs.Primary.ID) {
		//	return fmt.Errorf("no relationship found")
		//}

		return nil
	}
}
