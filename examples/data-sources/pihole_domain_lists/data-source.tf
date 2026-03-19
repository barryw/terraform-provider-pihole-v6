data "pihole_domain_lists" "deny_exact" {
  type = "deny"
  kind = "exact"
}

output "deny_exact_domains" {
  value = data.pihole_domain_lists.deny_exact.domains
}
