resource "pihole_setting" "app_sudo" {
  key   = "webserver.api.app_sudo"
  value = jsonencode(true)
}

resource "pihole_setting" "block_ttl" {
  key   = "dns.blockTTL"
  value = jsonencode(2)
}
