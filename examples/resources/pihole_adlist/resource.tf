resource "pihole_adlist" "stevenblack" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
  type    = "block"
  comment = "StevenBlack unified hosts"
  groups  = [0]
  enabled = true
}
