resource "pihole_client" "desktop" {
  client  = "192.168.1.100"
  comment = "Barry's desktop"
  groups  = [0]
}
