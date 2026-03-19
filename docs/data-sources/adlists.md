---
page_title: "pihole_adlists Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all adlists from Pi-hole.
---

# pihole_adlists (Data Source)

Fetches all adlists configured in Pi-hole.

## Example Usage

```terraform
data "pihole_adlists" "all" {}

output "all_adlists" {
  value = data.pihole_adlists.all.adlists
}
```

## Schema

### Read-Only

- `adlists` (List of Object) - List of all adlists. Each adlist contains:
  - `id` (String) - The address of the adlist.
  - `address` (String) - The URL of the adlist.
  - `type` (String) - The type of the adlist: `block` or `allow`.
  - `comment` (String) - The comment for the adlist.
  - `groups` (List of Number) - List of group IDs this adlist is assigned to.
  - `enabled` (Boolean) - Whether the adlist is enabled.
