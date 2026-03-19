---
page_title: "pihole_cname_record Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages a CNAME record in Pi-hole.
---

# pihole_cname_record (Resource)

Manages a CNAME record in Pi-hole's local DNS. This resource creates an alias from one domain to another.

CNAME records are immutable -- changing `domain`, `target`, or `ttl` forces replacement of the resource.

## Example Usage

```terraform
# Point files.home.lan to nas.home.lan
resource "pihole_cname_record" "files" {
  domain = "files.home.lan"
  target = "nas.home.lan"
}

# CNAME with explicit TTL
resource "pihole_cname_record" "media" {
  domain = "media.home.lan"
  target = "nas.home.lan"
  ttl    = 300
}
```

## Schema

### Required

- `domain` (String) - The domain name for the CNAME record. Changing this forces a new resource.
- `target` (String) - The target hostname for the CNAME record. Changing this forces a new resource.

### Optional

- `ttl` (Number) - Time to live in seconds. `0` means use the default. Changing this forces a new resource.

### Read-Only

- `id` (String) - Composite ID in the format `domain:target`.

## Import

Import is supported using the following syntax:

```shell
# Import format: domain:target
terraform import pihole_cname_record.files "files.home.lan:nas.home.lan"
```
