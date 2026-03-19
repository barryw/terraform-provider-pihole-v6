---
page_title: "pihole_adlist Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages an adlist in Pi-hole.
---

# pihole_adlist (Resource)

Manages an adlist (block or allow list) in Pi-hole. Adlists are URLs pointing to lists of domains that Pi-hole will use for blocking or allowing queries.

The `address` attribute is immutable -- changing it forces replacement. The `type`, `comment`, `groups`, and `enabled` attributes can be updated in-place.

## Example Usage

```terraform
# Block list
resource "pihole_adlist" "steven_black" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
  type    = "block"
  comment = "Steven Black unified hosts list"
}

# Allow list assigned to a specific group
resource "pihole_adlist" "allowlist" {
  address = "https://example.com/allowlist.txt"
  type    = "allow"
  comment = "Curated allow list"
  groups  = [0, 1]
  enabled = true
}
```

## Schema

### Required

- `address` (String) - The URL of the adlist. Changing this forces a new resource.
- `type` (String) - The type of the adlist: `block` or `allow`.

### Optional

- `comment` (String) - An optional comment for the adlist.
- `groups` (List of Number) - List of group IDs this adlist is assigned to.
- `enabled` (Boolean) - Whether the adlist is enabled. Defaults to `true`.

### Read-Only

- `id` (String) - The address of the adlist (same as `address`).

## Import

Import is supported using the following syntax:

```shell
# Import format: the adlist URL
terraform import pihole_adlist.steven_black "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
```
