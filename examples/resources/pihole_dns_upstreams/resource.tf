# Manage the ordered list of upstream DNS servers. Singleton — define at most once.
resource "pihole_dns_upstreams" "main" {
  upstreams = [
    "1.1.1.1",
    "1.0.0.1",
    "127.0.0.1#5335", # e.g. a local unbound resolver
  ]
}
