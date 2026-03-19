data "pihole_dns_records" "all" {}

output "all_dns_records" {
  value = data.pihole_dns_records.all.records
}
