---
page_title: "pihole_groups Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all Pi-hole groups.
---

# pihole_groups (Data Source)

Lists all groups configured in Pi-hole.

## Example Usage

```terraform
data "pihole_groups" "all" {}

output "all_groups" {
  value = data.pihole_groups.all.groups
}
```

## Schema

### Read-Only

- `groups` (List of Object) - List of all groups. Each group contains:
  - `id` (String) - The group name (same as `name`).
  - `name` (String) - The name of the group.
  - `comment` (String) - The comment associated with the group.
  - `enabled` (Boolean) - Whether the group is enabled.
