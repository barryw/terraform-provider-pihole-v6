---
page_title: "pihole_group Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single Pi-hole group by name.
---

# pihole_group (Data Source)

Fetches a single Pi-hole group by its name.

## Example Usage

```terraform
data "pihole_group" "default" {
  name = "Default"
}

output "default_group_enabled" {
  value = data.pihole_group.default.enabled
}
```

## Schema

### Required

- `name` (String) - The name of the group to look up.

### Read-Only

- `id` (String) - The group name (same as `name`).
- `comment` (String) - The comment associated with the group.
- `enabled` (Boolean) - Whether the group is enabled.
