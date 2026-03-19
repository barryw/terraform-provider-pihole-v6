---
page_title: "pihole_adlist Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single adlist by address.
---

# pihole_adlist (Data Source)

Fetches a single adlist from Pi-hole by its URL address.

## Example Usage

```terraform
data "pihole_adlist" "steven_black" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
}

output "steven_black_type" {
  value = data.pihole_adlist.steven_black.type
}

output "steven_black_enabled" {
  value = data.pihole_adlist.steven_black.enabled
}
```

## Schema

### Required

- `address` (String) - The URL of the adlist to look up.

### Read-Only

- `id` (String) - The address of the adlist (same as `address`).
- `type` (String) - The type of the adlist: `block` or `allow`.
- `comment` (String) - The comment for the adlist.
- `groups` (List of Number) - List of group IDs this adlist is assigned to.
- `enabled` (Boolean) - Whether the adlist is enabled.
