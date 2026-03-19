---
page_title: "pihole_domain_list Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches a single domain list entry by domain, type, and kind.
---

# pihole_domain_list (Data Source)

Fetches a single domain list entry from Pi-hole by its domain, type, and kind.

## Example Usage

```terraform
data "pihole_domain_list" "blocked" {
  domain = "ads.example.com"
  type   = "deny"
  kind   = "exact"
}

output "blocked_enabled" {
  value = data.pihole_domain_list.blocked.enabled
}

output "blocked_groups" {
  value = data.pihole_domain_list.blocked.groups
}
```

## Schema

### Required

- `domain` (String) - The domain or regex pattern.
- `type` (String) - The list type: `allow` or `deny`.
- `kind` (String) - The match kind: `exact` or `regex`.

### Read-Only

- `id` (String) - Composite ID in the format `type:kind:domain`.
- `comment` (String) - Comment for this domain entry.
- `groups` (List of Number) - List of group IDs this domain entry belongs to.
- `enabled` (Boolean) - Whether this domain entry is enabled.
