---
page_title: "pihole_cname_records Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all CNAME records from Pi-hole.
---

# pihole_cname_records (Data Source)

Fetches all CNAME records from Pi-hole.

## Example Usage

```terraform
data "pihole_cname_records" "all" {}

output "all_cname_records" {
  value = data.pihole_cname_records.all.records
}
```

## Schema

### Read-Only

- `records` (List of Object) - List of all CNAME records. Each record contains:
  - `id` (String) - Composite ID in the format `domain:target`.
  - `domain` (String) - The domain name for the CNAME record.
  - `target` (String) - The target hostname for the CNAME record.
  - `ttl` (Number) - Time to live in seconds.
