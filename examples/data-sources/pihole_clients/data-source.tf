data "pihole_clients" "all" {}

output "all_clients" {
  value = data.pihole_clients.all.clients
}
