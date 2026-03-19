# terraform-provider-pihole-v6

[![CI](https://github.com/barryw/terraform-provider-pihole-v6/actions/workflows/ci.yml/badge.svg)](https://github.com/barryw/terraform-provider-pihole-v6/actions/workflows/ci.yml)
[![Release](https://github.com/barryw/terraform-provider-pihole-v6/actions/workflows/release.yml/badge.svg)](https://github.com/barryw/terraform-provider-pihole-v6/actions/workflows/release.yml)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)

A Terraform/OpenTofu provider for managing [Pi-hole](https://pi-hole.net/) v6 configuration through its API. Built from scratch on the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).

This provider targets **Pi-hole v6 exclusively**. It is not compatible with Pi-hole v5 or earlier.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.6 or [OpenTofu](https://opentofu.org/) >= 1.6
- Pi-hole v6 with an [app-password](https://docs.pi-hole.net/api/auth/) configured
- Go >= 1.25 (only for building from source)

## Installation

Add the provider to your `required_providers` block:

```hcl
terraform {
  required_providers {
    pihole = {
      source = "barryw/pihole-v6"
    }
  }
}
```

## Provider Configuration

The provider requires a Pi-hole URL and app-password. These can be set directly in the provider block or via environment variables.

```hcl
provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}
```

### Environment Variables

| Variable | Description |
|---|---|
| `PIHOLE_URL` | Base URL of the Pi-hole instance (e.g. `http://192.168.1.1:8080`) |
| `PIHOLE_PASSWORD` | App-password for the Pi-hole API |

Environment variables are used as fallbacks when the corresponding provider attribute is not set.

### Multiple Instances

Use provider aliases to manage multiple Pi-hole instances:

```hcl
provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password_primary
}

provider "pihole" {
  alias    = "secondary"
  url      = "http://192.168.1.2:8080"
  password = var.pihole_password_secondary
}
```

## Quick Start

This example manages DNS records on two Pi-hole instances:

```hcl
terraform {
  required_providers {
    pihole = {
      source = "barryw/pihole-v6"
    }
  }
}

variable "pihole_password" {
  type      = string
  sensitive = true
}

provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}

provider "pihole" {
  alias    = "secondary"
  url      = "http://192.168.1.2:8080"
  password = var.pihole_password
}

# Create a group for IoT devices
resource "pihole_group" "iot" {
  name    = "IoT Devices"
  comment = "Group for IoT device filtering"
}

# Add a local DNS record on the primary instance
resource "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
  ip     = "192.168.1.100"
}

# Add the same record on the secondary instance
resource "pihole_dns_record" "nas_secondary" {
  provider = pihole.secondary
  domain   = "nas.home.lan"
  ip       = "192.168.1.100"
}

# Add a CNAME alias
resource "pihole_cname_record" "files" {
  domain = "files.home.lan"
  target = "nas.home.lan"
}

# Block a domain
resource "pihole_domain_list" "block_ads" {
  domain = "ads.example.com"
  type   = "deny"
  kind   = "exact"
}

# Add an adlist
resource "pihole_adlist" "steven_black" {
  address = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
  type    = "block"
  comment = "Steven Black unified hosts"
}

# Register a client
resource "pihole_client" "living_room_tv" {
  client  = "192.168.1.50"
  comment = "Living room smart TV"
  groups  = [0]
}
```

## Resources

| Resource | Description |
|---|---|
| [pihole_dns_record](docs/resources/dns_record.md) | Manages a local DNS A/AAAA record |
| [pihole_cname_record](docs/resources/cname_record.md) | Manages a local CNAME record |
| [pihole_group](docs/resources/group.md) | Manages a Pi-hole group |
| [pihole_adlist](docs/resources/adlist.md) | Manages an adlist (block or allow list URL) |
| [pihole_domain_list](docs/resources/domain_list.md) | Manages a domain list entry (allow/deny, exact/regex) |
| [pihole_client](docs/resources/client.md) | Manages a client definition |

## Data Sources

| Data Source | Description |
|---|---|
| [pihole_dns_record](docs/data-sources/dns_record.md) | Fetches a single DNS record by domain |
| [pihole_dns_records](docs/data-sources/dns_records.md) | Fetches all DNS records |
| [pihole_cname_record](docs/data-sources/cname_record.md) | Fetches a single CNAME record by domain |
| [pihole_cname_records](docs/data-sources/cname_records.md) | Fetches all CNAME records |
| [pihole_group](docs/data-sources/group.md) | Fetches a single group by name |
| [pihole_groups](docs/data-sources/groups.md) | Fetches all groups |
| [pihole_adlist](docs/data-sources/adlist.md) | Fetches a single adlist by address |
| [pihole_adlists](docs/data-sources/adlists.md) | Fetches all adlists |
| [pihole_domain_list](docs/data-sources/domain_list.md) | Fetches a single domain list entry |
| [pihole_domain_lists](docs/data-sources/domain_lists.md) | Fetches all domain list entries (with optional filters) |
| [pihole_client](docs/data-sources/client.md) | Fetches a single client by identifier |
| [pihole_clients](docs/data-sources/clients.md) | Fetches all clients |

## Import

All resources support `terraform import`. The import ID format varies by resource type:

```shell
# DNS record: domain:ip
terraform import pihole_dns_record.example "nas.home.lan:192.168.1.100"

# CNAME record: domain:target
terraform import pihole_cname_record.example "files.home.lan:nas.home.lan"

# Group: group name
terraform import pihole_group.example "IoT Devices"

# Adlist: the adlist URL
terraform import pihole_adlist.example "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"

# Domain list entry: type:kind:domain
terraform import pihole_domain_list.example "deny:exact:ads.example.com"

# Client: the client identifier
terraform import pihole_client.example "192.168.1.50"
```

## Development

### Building from Source

```shell
git clone https://github.com/barryw/terraform-provider-pihole-v6.git
cd terraform-provider-pihole-v6
go build -o terraform-provider-pihole-v6
```

### Running Tests

```shell
go test ./...
```

### Acceptance Tests

Acceptance tests run against a real Pi-hole v6 instance. Set the required environment variables before running:

```shell
export PIHOLE_URL="http://localhost:8080"
export PIHOLE_PASSWORD="your-app-password"
TF_ACC=1 go test ./internal/provider/ -v
```

## License

[Mozilla Public License 2.0](LICENSE)
