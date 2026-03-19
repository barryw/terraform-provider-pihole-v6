data "pihole_cname_records" "all" {}

output "all_cname_records" {
  value = data.pihole_cname_records.all.records
}
