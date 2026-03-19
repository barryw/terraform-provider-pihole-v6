package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecordsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "pihole_dns_record" "test" {
  domain = "acc-ds-test.example"
  ip     = "10.99.99.2"
}

data "pihole_dns_records" "all" {
  depends_on = [pihole_dns_record.test]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pihole_dns_records.all", "records.#"),
				),
			},
		},
	})
}
