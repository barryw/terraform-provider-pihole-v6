resource "pihole_group" "iot" {
  name    = "IoT Devices"
  comment = "Group for IoT devices with restricted access"
  enabled = true
}
