data "pihole_group" "iot" {
  name = "IoT Devices"
}

output "iot_group_enabled" {
  value = data.pihole_group.iot.enabled
}
