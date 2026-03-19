terraform {
  required_providers {
    pihole = {
      source = "barryw/pihole-v6"
    }
  }
}

# Configure the PiHole provider
provider "pihole" {
  url      = "http://192.168.1.1:8080"
  password = var.pihole_password
}

# Multiple instances using aliases
provider "pihole" {
  alias    = "secondary"
  url      = "http://192.168.1.2:8080"
  password = var.pihole_password
}

variable "pihole_password" {
  type      = string
  sensitive = true
}
