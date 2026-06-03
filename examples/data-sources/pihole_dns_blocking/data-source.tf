# Read the current DNS blocking state.
data "pihole_dns_blocking" "current" {}

output "blocking_enabled" {
  value = data.pihole_dns_blocking.current.enabled
}
