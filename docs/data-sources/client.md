---
page_title: "pihole_client Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single client by its identifier.
---

# pihole_client (Data Source)

Fetches a single client from Pi-hole by its identifier (IP address, MAC address, CIDR range, or hostname).

## Example Usage

```terraform
data "pihole_client" "tv" {
  client = "192.168.1.50"
}

output "tv_comment" {
  value = data.pihole_client.tv.comment
}

output "tv_groups" {
  value = data.pihole_client.tv.groups
}
```

## Schema

### Required

- `client` (String) - The client identifier (IP address, MAC address, CIDR range, or hostname) to look up.

### Read-Only

- `id` (String) - The client identifier (same as `client`).
- `comment` (String) - Comment for the client.
- `groups` (List of Number) - List of group IDs the client belongs to.
