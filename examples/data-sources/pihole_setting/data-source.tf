data "pihole_setting" "block_ttl" {
  key = "dns.blockTTL"
}

output "block_ttl" {
  value = jsondecode(data.pihole_setting.block_ttl.value)
}
