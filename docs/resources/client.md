---
page_title: "pihole_client Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages a client in Pi-hole.
---

# pihole_client (Resource)

Manages a client definition in Pi-hole. Clients can be identified by IP address, MAC address, CIDR range, or hostname, and can be assigned to groups for differentiated filtering policies.

The `client` attribute is immutable -- changing it forces replacement. The `comment` and `groups` attributes can be updated in-place.

## Example Usage

```terraform
# Client by IP address
resource "pihole_client" "living_room_tv" {
  client  = "192.168.1.50"
  comment = "Living room smart TV"
  groups  = [0]
}

# Client by MAC address
resource "pihole_client" "laptop" {
  client  = "AA:BB:CC:DD:EE:FF"
  comment = "Work laptop"
  groups  = [0, 1]
}

# Client by CIDR range
resource "pihole_client" "iot_subnet" {
  client  = "192.168.10.0/24"
  comment = "IoT VLAN"
  groups  = [2]
}
```

## Schema

### Required

- `client` (String) - The client identifier (IP address, MAC address, CIDR range, or hostname). Changing this forces a new resource.

### Optional

- `comment` (String) - Optional comment for the client.
- `groups` (List of Number) - List of group IDs the client belongs to.

### Read-Only

- `id` (String) - The client identifier (same as `client`).

## Import

Import is supported using the following syntax:

```shell
# Import format: the client identifier
terraform import pihole_client.living_room_tv "192.168.1.50"
```
