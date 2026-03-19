---
page_title: "pihole_group Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages a Pi-hole group.
---

# pihole_group (Resource)

Manages a group in Pi-hole. Groups are used to organize clients, adlists, and domain list entries so that different filtering rules can be applied to different sets of clients.

Groups support in-place updates for `name`, `comment`, and `enabled`.

## Example Usage

```terraform
# Group for IoT devices with stricter filtering
resource "pihole_group" "iot" {
  name    = "IoT Devices"
  comment = "Strict filtering for IoT devices"
}

# Disabled group for testing
resource "pihole_group" "testing" {
  name    = "Testing"
  comment = "Temporarily disabled group"
  enabled = false
}
```

## Schema

### Required

- `name` (String) - The name of the group.

### Optional

- `comment` (String) - An optional comment for the group. Defaults to `""`.
- `enabled` (Boolean) - Whether the group is enabled. Defaults to `true`.

### Read-Only

- `id` (String) - The group name (same as `name`).

## Import

Import is supported using the following syntax:

```shell
# Import format: group name
terraform import pihole_group.iot "IoT Devices"
```
