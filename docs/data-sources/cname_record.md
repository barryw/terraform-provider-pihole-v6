---
page_title: "pihole_cname_record Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single CNAME record by domain name.
---

# pihole_cname_record (Data Source)

Fetches a single CNAME record from Pi-hole by domain name.

## Example Usage

```terraform
data "pihole_cname_record" "files" {
  domain = "files.home.lan"
}

output "files_target" {
  value = data.pihole_cname_record.files.target
}

output "files_ttl" {
  value = data.pihole_cname_record.files.ttl
}
```

## Schema

### Required

- `domain` (String) - The domain name to look up.

### Read-Only

- `id` (String) - Composite ID in the format `domain:target`.
- `target` (String) - The target hostname for the CNAME record.
- `ttl` (Number) - Time to live in seconds.
