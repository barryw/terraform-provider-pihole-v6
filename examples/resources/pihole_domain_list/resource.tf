resource "pihole_domain_list" "block_ads" {
  domain  = "ads.example.com"
  type    = "deny"
  kind    = "exact"
  comment = "Block ads from example.com"
  groups  = [0]
  enabled = true
}
