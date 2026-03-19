---
page_title: "pihole_dns_record Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages a local DNS record in Pi-hole.
---

# pihole_dns_record (Resource)

Manages a local DNS A/AAAA record in Pi-hole. This resource maps a domain name to an IP address in Pi-hole's local DNS resolver.

DNS records are immutable -- changing either `domain` or `ip` forces replacement of the resource.

## Example Usage

```terraform
# A record for a NAS
resource "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
  ip     = "192.168.1.100"
}

# AAAA record
resource "pihole_dns_record" "nas_v6" {
  domain = "nas.home.lan"
  ip     = "fd00::100"
}

# Kubernetes ingress
resource "pihole_dns_record" "grafana" {
  domain = "grafana.home.lan"
  ip     = "192.168.1.200"
}
```

## Schema

### Required

- `domain` (String) - The domain name. Changing this forces a new resource.
- `ip` (String) - The IP address to resolve to. Changing this forces a new resource.

### Read-Only

- `id` (String) - Composite ID in the format `domain:ip`.

## Import

Import is supported using the following syntax:

```shell
# Import format: domain:ip
terraform import pihole_dns_record.nas "nas.home.lan:192.168.1.100"
```
