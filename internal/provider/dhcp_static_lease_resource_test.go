package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDHCPStaticLeaseResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and verify
			{
				Config: `
resource "pihole_dhcp_static_lease" "test" {
  mac      = "de:ad:be:ef:00:01"
  ip       = "10.99.99.50"
  hostname = "acc-test-host"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "mac", "de:ad:be:ef:00:01"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "ip", "10.99.99.50"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "hostname", "acc-test-host"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "id", "de:ad:be:ef:00:01:10.99.99.50"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_dhcp_static_lease.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
