# Manage PiHole's global DNS blocking state. Singleton — define at most once.
resource "pihole_dns_blocking" "main" {
  enabled = true
}
