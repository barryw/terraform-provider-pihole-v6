package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCNAMERecordResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "pihole_cname_record" "test" {
  domain = "acc-alias.example"
  target = "acc-target.example"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_cname_record.test", "domain", "acc-alias.example"),
					resource.TestCheckResourceAttr("pihole_cname_record.test", "target", "acc-target.example"),
					resource.TestCheckResourceAttr("pihole_cname_record.test", "ttl", "0"),
				),
			},
			{
				ResourceName:      "pihole_cname_record.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCNAMERecordResource_withTTL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "pihole_cname_record" "test_ttl" {
  domain = "acc-ttl.example"
  target = "acc-target.example"
  ttl    = 300
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_cname_record.test_ttl", "domain", "acc-ttl.example"),
					resource.TestCheckResourceAttr("pihole_cname_record.test_ttl", "target", "acc-target.example"),
					resource.TestCheckResourceAttr("pihole_cname_record.test_ttl", "ttl", "300"),
				),
			},
		},
	})
}
