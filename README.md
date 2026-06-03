[![Build Status](https://ci.barrywalker.io/api/badges/49/status.svg)](https://ci.barrywalker.io/repos/49)
[![Latest Release](https://img.shields.io/github/v/release/barryw/terraform-provider-pihole-v6)](https://github.com/barryw/terraform-provider-pihole-v6/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)

# Terraform Provider for Pi-hole v6

Manage your [Pi-hole](https://pi-hole.net/) v6 infrastructure as code. This provider communicates with Pi-hole's v6 REST API to manage DNS records, CNAME records, groups, adlists, domain allow/deny lists, client assignments, and configuration settings.

Built from scratch on the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) targeting **Pi-hole v6 exclusively**.

## Why this provider?

Existing Pi-hole Terraform providers target the legacy v5 PHP API or have broken CSRF handling with v6. This provider was purpose-built for Pi-hole v6's new Go-based API (pihole-FTL) with proper session management and full resource coverage.

## Features

- **10 resources** with full CRUD, import, and drift detection
- **17 data sources** (singular lookup + plural list for most resource types)
- **DNS blocking & upstreams** — manage Pi-hole's global blocking toggle and upstream resolvers declaratively
- **Multi-instance support** via provider aliases — manage multiple Pi-holes in a single config
- **Generic settings management** via `pihole_setting` — configure any Pi-hole setting by dot-notation path (now JSON-validated and canonicalized to avoid spurious diffs)
- **Idempotent operations** — safe to re-apply, handles already-existing records gracefully
- **Concurrency-safe session management** — authenticates once per provider instance, collapses parallel re-authentication into a single request, and retries transient errors with context-aware backoff

## Resources

| Resource | Description |
|---|---|
| `pihole_dns_record` | Local DNS A/AAAA records |
| `pihole_cname_record` | Local CNAME records with optional TTL |
| `pihole_group` | Groups for organizing clients and lists |
| `pihole_adlist` | Block/allow list URLs |
| `pihole_domain_list` | Domain allow/deny entries (exact or regex) |
| `pihole_client` | Client assignments by IP, MAC, or CIDR |
| `pihole_dhcp_static_lease` | DHCP static lease reservations |
| `pihole_setting` | Any Pi-hole configuration setting by path |
| `pihole_dns_blocking` | Global DNS blocking on/off state (singleton) |
| `pihole_dns_upstreams` | Ordered list of upstream DNS servers (singleton) |

## Data Sources

Most resources have both a singular (lookup by key) and plural (list all) data source. The singleton settings (`dns_blocking`, `dns_upstreams`) expose a single data source each:

| Data Source | Description |
|---|---|
| `pihole_dns_record` / `pihole_dns_records` | Read DNS records |
| `pihole_cname_record` / `pihole_cname_records` | Read CNAME records |
| `pihole_group` / `pihole_groups` | Read groups |
| `pihole_adlist` / `pihole_adlists` | Read adlists |
| `pihole_domain_list` / `pihole_domain_lists` | Read domain entries (with optional type/kind filters) |
| `pihole_client` / `pihole_clients` | Read client assignments |
| `pihole_dhcp_static_lease` / `pihole_dhcp_static_leases` | Read DHCP static leases |
| `pihole_setting` | Read any configuration setting |
| `pihole_dns_blocking` | Read the current DNS blocking state |
| `pihole_dns_upstreams` | Read the configured upstream DNS servers |

## Quick Start

```hcl
terraform {
  required_providers {
    pihole = {
      source  = "barryw/pihole-v6"
      version = "~> 0.1"
    }
  }
}

provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}

# Manage DNS records across your network
resource "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
  ip     = "192.168.1.100"
}

resource "pihole_cname_record" "files" {
  domain = "files.home.lan"
  target = "nas.home.lan"
}

# Organize devices into groups
resource "pihole_group" "iot" {
  name    = "IoT Devices"
  comment = "Restricted DNS for IoT"
}

# Manage blocklists
resource "pihole_adlist" "steven_black" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
  type    = "block"
  comment = "Steven Black unified hosts"
  groups  = [0]
}

# Configure Pi-hole settings
resource "pihole_setting" "cache_size" {
  key   = "dns.cache.size"
  value = jsonencode(20000)
}
```

## Multi-Instance Management

Use provider aliases to keep multiple Pi-hole instances in sync:

```hcl
provider "pihole" {
  alias    = "primary"
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}

provider "pihole" {
  alias    = "secondary"
  url      = "http://192.168.1.2:8080"
  password = var.pihole_password
}

locals {
  dns_records = {
    "nas.lan"     = "192.168.1.100"
    "printer.lan" = "192.168.1.101"
    "camera.lan"  = "192.168.1.102"
  }
}

resource "pihole_dns_record" "primary" {
  for_each = local.dns_records
  provider = pihole.primary
  domain   = each.key
  ip       = each.value
}

resource "pihole_dns_record" "secondary" {
  for_each = local.dns_records
  provider = pihole.secondary
  domain   = each.key
  ip       = each.value
}
```

## Prerequisites

Before using this provider, each Pi-hole instance needs:

1. **An app-password** — generate one in Pi-hole web UI under Settings > API
2. **`app_sudo` enabled** — required for write operations

```bash
# Enable via CLI
sudo pihole-FTL --config webserver.api.app_sudo true

# Or via pihole.toml
# [webserver.api]
#   app_sudo = true
```

> **Note:** Without `app_sudo`, all write operations fail with `403 Forbidden`. This is separate from `allow_destructive`.

## Import

All resources support `terraform import`:

```shell
terraform import pihole_dns_record.example "nas.lan:192.168.1.100"
terraform import pihole_cname_record.example "files.lan:nas.lan"
terraform import pihole_group.example "IoT Devices"
terraform import pihole_adlist.example "https://example.com/blocklist.txt"
terraform import pihole_domain_list.example "deny:exact:ads.example.com"
terraform import pihole_client.example "192.168.1.50"
terraform import pihole_dhcp_static_lease.example "aa:bb:cc:dd:ee:ff:192.168.1.50"
terraform import pihole_setting.example "dns.cache.size"

# Singletons — import with their fixed id
terraform import pihole_dns_blocking.example "blocking"
terraform import pihole_dns_upstreams.example "upstreams"
```

## Environment Variables

| Variable | Description |
|---|---|
| `PIHOLE_URL` | Base URL (fallback when `url` is not set) |
| `PIHOLE_PASSWORD` | App-password (fallback when `password` is not set) |

## Development

```shell
# Build
git clone https://github.com/barryw/terraform-provider-pihole-v6.git
cd terraform-provider-pihole-v6
go build -o terraform-provider-pihole-v6

# Unit tests
go test ./...

# Acceptance tests (requires a running Pi-hole v6)
docker compose up -d --wait
PIHOLE_URL=http://localhost:18080 PIHOLE_PASSWORD=test-password TF_ACC=1 go test ./internal/provider/ -v

# Regenerate docs
make docs
```

### Generating documentation

`make docs` runs [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs), which **requires a real `terraform` CLI on your `PATH`** to introspect the provider schema. **OpenTofu (`tofu`) does not work** for this step: `tfplugindocs` addresses the provider under `registry.terraform.io/hashicorp/…`, while `tofu` resolves the same source under `registry.opentofu.org/…`, so `terraform init` can't find the locally-built provider. Install a `terraform` binary (1.5.7 is the last MPL-licensed release) — it's only used for doc generation, not at runtime.

## License

[Mozilla Public License 2.0](LICENSE)

## Related

- [go-pihole](https://github.com/barryw/go-pihole) — standalone Go client library for the Pi-hole v6 API (used by this provider)

