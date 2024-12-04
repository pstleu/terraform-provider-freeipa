package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeipaHostDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFreeipaHostDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_host.test", "id", "duba-nfws-sgwa01.corp.example.com"),
					resource.TestCheckResourceAttr("data.freeipa_host.test", "fqdn", "duba-nfws-sgwa01.corp.example.com"),
					resource.TestCheckResourceAttr("data.freeipa_host.test", "hostname", "duba-nfws-sgwa01"),
				),
			},
		},
	})
}

const testAccFreeipaHostDataSourceConfig = `
data "freeipa_host" "test" {
  fqdn = "duba-nfws-sgwa01.corp.example.com"
}

provider "freeipa" {
	host     = "duba-shp-doma01.corp.example.com"
	username = "terraform"
	password = "test"
	realm    = "CORP.EXAMPLE.COM"
	insecure = true
}

`
