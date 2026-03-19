data "pihole_dns_record" "myhost" {
  domain = "myhost.lan"
}

output "myhost_ip" {
  value = data.pihole_dns_record.myhost.ip
}
