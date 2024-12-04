# Copyright (c) HashiCorp, Inc.

data "freeipa_host" "example" {
  fqdn = "duba-nfws-otfa01.corp.example.com"
}

output "fqdn" {
  value = data.freeipa_host.example
}
