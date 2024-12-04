# Copyright (c) HashiCorp, Inc.

data "freeipa_host" "example" {
  fqdn = "duba-nfws-otfa01.corp.trimbletl.com"
}

output "fqdn" {
  value = data.freeipa_host.example
}

resource "freeipa_host" "example" {
  fqdn        = "ipa-tf-provider-test2.corp.trimbletl.com"
  force       = true
  description = "Test host created by Terraform 4"
}
