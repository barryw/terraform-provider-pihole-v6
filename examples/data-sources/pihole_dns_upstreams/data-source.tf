# Read the currently configured upstream DNS servers.
data "pihole_dns_upstreams" "current" {}

output "upstreams" {
  value = data.pihole_dns_upstreams.current.upstreams
}
