package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdlistResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_adlist" "test" {
  address = "https://acc-test.example/hosts.txt"
  type    = "block"
  comment = "acceptance test adlist"
  groups  = [0]
  enabled = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_adlist.test", "address", "https://acc-test.example/hosts.txt"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "type", "block"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "comment", "acceptance test adlist"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "id", "https://acc-test.example/hosts.txt"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "groups.0", "0"),
				),
			},
			// Update comment and groups
			{
				Config: `
resource "pihole_adlist" "test" {
  address = "https://acc-test.example/hosts.txt"
  type    = "block"
  comment = "updated adlist comment"
  groups  = [0]
  enabled = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_adlist.test", "comment", "updated adlist comment"),
					resource.TestCheckResourceAttr("pihole_adlist.test", "enabled", "false"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_adlist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
