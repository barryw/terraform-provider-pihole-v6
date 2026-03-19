---
page_title: "pihole_dns_records Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all DNS records from Pi-hole.
---

# pihole_dns_records (Data Source)

Fetches all local DNS records from Pi-hole.

## Example Usage

```terraform
data "pihole_dns_records" "all" {}

output "all_dns_records" {
  value = data.pihole_dns_records.all.records
}
```

## Schema

### Read-Only

- `records` (List of Object) - List of all DNS records. Each record contains:
  - `id` (String) - Composite ID in the format `domain:ip`.
  - `domain` (String) - The domain name.
  - `ip` (String) - The IP address the domain resolves to.
