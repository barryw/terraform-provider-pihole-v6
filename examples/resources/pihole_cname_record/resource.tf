resource "pihole_cname_record" "example" {
  domain = "alias.lan"
  target = "myhost.lan"
  ttl    = 300
}
