package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecordResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and verify
			{
				Config: `
resource "pihole_dns_record" "test" {
  domain = "acc-test.example"
  ip     = "10.99.99.1"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_record.test", "domain", "acc-test.example"),
					resource.TestCheckResourceAttr("pihole_dns_record.test", "ip", "10.99.99.1"),
					resource.TestCheckResourceAttr("pihole_dns_record.test", "id", "acc-test.example:10.99.99.1"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_dns_record.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
