# Terraform Freeipa Provider using Terraform Plugin Framework

_This repository is built from terraform provider scafffolding repo using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) SDK. _


- Freeipa Host resource and a data source can be found `internal/provider/`,
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```
## Installing Dependencies

This provider uses [Go Freeipa SDK](https://github.com/camptocamp/go-freeipa) for interacting with Freeipa server.

```shell
go get github.com/ccin2p3/go-freeipa
go mod tidy
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```terraform
terraform {
  required_providers {
    freeipa = {
      source = "hashicorp.com/mashanm/freeipa"
    }
  }
}

provider "freeipa" {
  host     = "duba-shp-doma01.corp.example.com"
  username = "terraform"
  password = "password"
  realm = "corp.example.com"
}

data "freeipa_host" "example" {
  fqdn = "duba-nfws-otfa01.corp.example.com"
}

output "fqdn" {
  value = data.freeipa_host.example
}

```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
