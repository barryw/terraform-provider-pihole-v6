---
page_title: "pihole_clients Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all clients from Pi-hole.
---

# pihole_clients (Data Source)

Fetches all client definitions from Pi-hole.

## Example Usage

```terraform
data "pihole_clients" "all" {}

output "all_clients" {
  value = data.pihole_clients.all.clients
}
```

## Schema

### Read-Only

- `clients` (List of Object) - List of all clients. Each client contains:
  - `id` (String) - The client identifier.
  - `client` (String) - The client identifier (IP address, MAC address, CIDR range, or hostname).
  - `comment` (String) - Comment for the client.
  - `groups` (List of Number) - List of group IDs the client belongs to.
