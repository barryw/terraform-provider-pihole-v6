package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainListResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: `
resource "pihole_domain_list" "test" {
  domain  = "acc-blocked.example"
  type    = "deny"
  kind    = "exact"
  comment = "acceptance test domain"
  enabled = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_domain_list.test", "domain", "acc-blocked.example"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "type", "deny"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "kind", "exact"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "comment", "acceptance test domain"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "id", "deny:exact:acc-blocked.example"),
				),
			},
			// Update comment and enabled
			{
				Config: `
resource "pihole_domain_list" "test" {
  domain  = "acc-blocked.example"
  type    = "deny"
  kind    = "exact"
  comment = "updated domain comment"
  enabled = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_domain_list.test", "comment", "updated domain comment"),
					resource.TestCheckResourceAttr("pihole_domain_list.test", "enabled", "false"),
				),
			},
			// Import
			{
				ResourceName:      "pihole_domain_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDomainListResource_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "pihole_domain_list" "test_regex" {
  domain = "^acc-.*\\.example$"
  type   = "deny"
  kind   = "regex"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_domain_list.test_regex", "domain", "^acc-.*\\.example$"),
					resource.TestCheckResourceAttr("pihole_domain_list.test_regex", "type", "deny"),
					resource.TestCheckResourceAttr("pihole_domain_list.test_regex", "kind", "regex"),
					resource.TestCheckResourceAttr("pihole_domain_list.test_regex", "enabled", "true"),
				),
			},
		},
	})
}
