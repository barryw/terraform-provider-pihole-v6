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
			// Update IP and hostname
			{
				Config: `
resource "pihole_dhcp_static_lease" "test" {
  mac      = "de:ad:be:ef:00:01"
  ip       = "10.99.99.51"
  hostname = "acc-test-host-updated"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "mac", "de:ad:be:ef:00:01"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "ip", "10.99.99.51"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "hostname", "acc-test-host-updated"),
					resource.TestCheckResourceAttr("pihole_dhcp_static_lease.test", "id", "de:ad:be:ef:00:01:10.99.99.51"),
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

func TestAccDHCPStaticLeasesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "pihole_dhcp_static_lease" "test" {
  mac      = "de:ad:be:ef:00:02"
  ip       = "10.99.99.60"
  hostname = "acc-ds-test-host"
}

data "pihole_dhcp_static_leases" "all" {
  depends_on = [pihole_dhcp_static_lease.test]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pihole_dhcp_static_leases.all", "leases.#"),
				),
			},
		},
	})
}
