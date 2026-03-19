data "pihole_client" "desktop" {
  client = "192.168.1.100"
}

output "desktop_groups" {
  value = data.pihole_client.desktop.groups
}
