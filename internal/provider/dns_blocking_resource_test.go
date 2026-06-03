package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSBlockingResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Disable blocking
			{
				Config: `
resource "pihole_dns_blocking" "test" {
  enabled = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_blocking.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pihole_dns_blocking.test", "id", "blocking"),
				),
			},
			// Re-enable blocking
			{
				Config: `
resource "pihole_dns_blocking" "test" {
  enabled = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_blocking.test", "enabled", "true"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_dns_blocking.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
