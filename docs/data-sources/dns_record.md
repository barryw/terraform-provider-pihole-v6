---
page_title: "pihole_dns_record Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single DNS record by domain name.
---

# pihole_dns_record (Data Source)

Fetches a single local DNS record from Pi-hole by domain name.

## Example Usage

```terraform
data "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
}

output "nas_ip" {
  value = data.pihole_dns_record.nas.ip
}
```

## Schema

### Required

- `domain` (String) - The domain name to look up.

### Read-Only

- `id` (String) - Composite ID in the format `domain:ip`.
- `ip` (String) - The IP address the domain resolves to.
