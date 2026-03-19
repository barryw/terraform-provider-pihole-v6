resource "pihole_dns_record" "example" {
  domain = "myhost.lan"
  ip     = "192.168.1.100"
}
