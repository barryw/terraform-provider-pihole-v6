data "pihole_cname_record" "alias" {
  domain = "alias.lan"
}

output "alias_target" {
  value = data.pihole_cname_record.alias.target
}
