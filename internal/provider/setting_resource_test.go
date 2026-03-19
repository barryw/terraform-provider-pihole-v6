package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_setting" "test" {
  key   = "dns.blockTTL"
  value = jsonencode(5)
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_setting.test", "key", "dns.blockTTL"),
					resource.TestCheckResourceAttr("pihole_setting.test", "value", "5"),
					resource.TestCheckResourceAttr("pihole_setting.test", "id", "dns.blockTTL"),
				),
			},
			// Update value
			{
				Config: `
resource "pihole_setting" "test" {
  key   = "dns.blockTTL"
  value = jsonencode(10)
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_setting.test", "value", "10"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_setting.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
