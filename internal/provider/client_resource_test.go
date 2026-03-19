package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClientResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_client" "test" {
  client  = "10.99.99.100"
  comment = "acceptance test client"
  groups  = [0]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_client.test", "client", "10.99.99.100"),
					resource.TestCheckResourceAttr("pihole_client.test", "comment", "acceptance test client"),
					resource.TestCheckResourceAttr("pihole_client.test", "id", "10.99.99.100"),
					resource.TestCheckResourceAttr("pihole_client.test", "groups.#", "1"),
					resource.TestCheckResourceAttr("pihole_client.test", "groups.0", "0"),
				),
			},
			// Update comment
			{
				Config: `
resource "pihole_client" "test" {
  client  = "10.99.99.100"
  comment = "updated client comment"
  groups  = [0]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_client.test", "comment", "updated client comment"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_client.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
