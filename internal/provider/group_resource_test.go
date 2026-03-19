package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_group" "test" {
  name    = "acc-test-group"
  comment = "acceptance test group"
  enabled = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_group.test", "name", "acc-test-group"),
					resource.TestCheckResourceAttr("pihole_group.test", "comment", "acceptance test group"),
					resource.TestCheckResourceAttr("pihole_group.test", "enabled", "true"),
				),
			},
			// Update
			{
				Config: `
resource "pihole_group" "test" {
  name    = "acc-test-group"
  comment = "updated comment"
  enabled = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_group.test", "comment", "updated comment"),
					resource.TestCheckResourceAttr("pihole_group.test", "enabled", "false"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
