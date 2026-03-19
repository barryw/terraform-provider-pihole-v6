---
page_title: "pihole_domain_lists Data Source - terraform-provider-pihole-v6"
subcategory: ""
description: |-
  Fetches all domain list entries from Pi-hole, with optional filtering.
---

# pihole_domain_lists (Data Source)

Fetches all domain list entries from Pi-hole. Results can be filtered by `type` and/or `kind`.

## Example Usage

```terraform
# Fetch all domain list entries
data "pihole_domain_lists" "all" {}

# Fetch only exact deny entries
data "pihole_domain_lists" "exact_deny" {
  type = "deny"
  kind = "exact"
}

# Fetch only regex entries (both allow and deny)
data "pihole_domain_lists" "all_regex" {
  kind = "regex"
}

output "deny_exact_domains" {
  value = data.pihole_domain_lists.exact_deny.domains
}
```

## Schema

### Optional

- `type` (String) - Filter by list type: `allow` or `deny`.
- `kind` (String) - Filter by match kind: `exact` or `regex`.

### Read-Only

- `domains` (List of Object) - List of domain list entries. Each entry contains:
  - `id` (String) - Composite ID in the format `type:kind:domain`.
  - `domain` (String) - The domain or regex pattern.
  - `type` (String) - The list type: `allow` or `deny`.
  - `kind` (String) - The match kind: `exact` or `regex`.
  - `comment` (String) - Comment for this domain entry.
  - `groups` (List of Number) - List of group IDs this domain entry belongs to.
  - `enabled` (Boolean) - Whether this domain entry is enabled.
