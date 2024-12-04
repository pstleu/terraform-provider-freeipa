# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    freeipa = {
      source  = "registry.terraform.io/ttl/freeipa"
      version = "1.0.0"
    }
  }
}

provider "freeipa" {
  host     = "duba-shp-doma01.corp.trimbletl.com"
  username = "terraform"
  password = "password"
  realm    = "corp.trimbletl.com"
  insecure = true
}
