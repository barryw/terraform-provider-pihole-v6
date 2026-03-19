data "pihole_groups" "all" {}

output "all_groups" {
  value = data.pihole_groups.all.groups
}
