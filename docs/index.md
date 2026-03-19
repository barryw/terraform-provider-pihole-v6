---
page_title: "PiHole Provider"
subcategory: ""
description: |-
  The PiHole provider manages Pi-hole v6 configuration through its API.
---

# PiHole Provider

The PiHole provider allows you to manage [Pi-hole](https://pi-hole.net/) v6 configuration as infrastructure-as-code. It communicates with the Pi-hole v6 REST API to manage DNS records, CNAME records, groups, adlists, domain lists, and client definitions.

This provider targets Pi-hole v6 exclusively and requires an app-password for authentication.

## Example Usage

```terraform
terraform {
  required_providers {
    pihole = {
      source = "barryw/pihole-v6"
    }
  }
}

provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}

variable "pihole_password" {
  type      = string
  sensitive = true
}

resource "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
  ip     = "192.168.1.100"
}
```

### Multiple Instances

Use provider aliases to manage more than one Pi-hole instance:

```terraform
provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password_primary
}

provider "pihole" {
  alias    = "secondary"
  url      = "http://192.168.1.2:8080"
  password = var.pihole_password_secondary
}

resource "pihole_dns_record" "nas" {
  domain = "nas.home.lan"
  ip     = "192.168.1.100"
}

resource "pihole_dns_record" "nas_secondary" {
  provider = pihole.secondary
  domain   = "nas.home.lan"
  ip       = "192.168.1.100"
}
```

## Schema

### Optional

- `url` (String) - Base URL of the Pi-hole instance (e.g. `http://192.168.1.1:8080`). May also be set via the `PIHOLE_URL` environment variable.
- `password` (String, Sensitive) - App-password for the Pi-hole API. May also be set via the `PIHOLE_PASSWORD` environment variable.

### Environment Variables

| Variable | Description |
|---|---|
| `PIHOLE_URL` | Base URL of the Pi-hole instance. Used when `url` is not set in the provider block. |
| `PIHOLE_PASSWORD` | App-password for the Pi-hole API. Used when `password` is not set in the provider block. |
