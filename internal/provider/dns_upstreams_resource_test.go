package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSUpstreamsResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_dns_upstreams" "test" {
  upstreams = ["1.1.1.1", "1.0.0.1"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.#", "2"),
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.1", "1.0.0.1"),
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "id", "upstreams"),
				),
			},
			// Update (reorder + change)
			{
				Config: `
resource "pihole_dns_upstreams" "test" {
  upstreams = ["8.8.8.8", "8.8.4.4", "9.9.9.9"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.#", "3"),
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.0", "8.8.8.8"),
					resource.TestCheckResourceAttr("pihole_dns_upstreams.test", "upstreams.2", "9.9.9.9"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_dns_upstreams.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
