# instructing terraform to use digital ocean
terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Assigning the do_token to the token argument of the provider
provider "digitalocean" {
  token = var.do_token
}