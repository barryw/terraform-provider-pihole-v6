data "pihole_domain_list" "example" {
  domain = "ads.example.com"
  type   = "deny"
  kind   = "exact"
}

output "domain_list_enabled" {
  value = data.pihole_domain_list.example.enabled
}
