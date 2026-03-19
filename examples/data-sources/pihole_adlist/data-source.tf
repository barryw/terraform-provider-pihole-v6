data "pihole_adlist" "stevenblack" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
}

output "stevenblack_enabled" {
  value = data.pihole_adlist.stevenblack.enabled
}
