---
page_title: "pihole_domain_list Resource - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Manages a domain list entry in Pi-hole.
---

# pihole_domain_list (Resource)

Manages a domain list entry (allow/deny, exact/regex) in Pi-hole. Domain list entries are individual domain rules used for whitelisting or blacklisting specific domains or patterns.

The `domain`, `type`, and `kind` attributes are immutable -- changing any of them forces replacement. The `comment`, `groups`, and `enabled` attributes can be updated in-place.

## Example Usage

```terraform
# Exact deny (blacklist a specific domain)
resource "pihole_domain_list" "block_ads" {
  domain = "ads.example.com"
  type   = "deny"
  kind   = "exact"
  comment = "Block ads from example.com"
}

# Regex deny (block all subdomains matching a pattern)
resource "pihole_domain_list" "block_tracking" {
  domain  = "(^|\\.)tracking\\.example\\.com$"
  type    = "deny"
  kind    = "regex"
  comment = "Block tracking subdomains"
}

# Exact allow (whitelist a domain)
resource "pihole_domain_list" "allow_cdn" {
  domain  = "cdn.example.com"
  type    = "allow"
  kind    = "exact"
  comment = "Allow CDN for streaming"
  groups  = [0]
}
```

## Schema

### Required

- `domain` (String) - The domain or regex pattern. Changing this forces a new resource.
- `type` (String) - The list type: `allow` or `deny`. Changing this forces a new resource.
- `kind` (String) - The match kind: `exact` or `regex`. Changing this forces a new resource.

### Optional

- `comment` (String) - Optional comment for this domain entry.
- `groups` (List of Number) - List of group IDs this domain entry belongs to.
- `enabled` (Boolean) - Whether this domain entry is enabled. Defaults to `true`.

### Read-Only

- `id` (String) - Composite ID in the format `type:kind:domain`.

## Import

Import is supported using the following syntax:

```shell
# Import format: type:kind:domain
terraform import pihole_domain_list.block_ads "deny:exact:ads.example.com"
```
