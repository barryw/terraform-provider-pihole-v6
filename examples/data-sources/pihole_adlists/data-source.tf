data "pihole_adlists" "all" {}

output "all_adlists" {
  value = data.pihole_adlists.all.adlists
}
