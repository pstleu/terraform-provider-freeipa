package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.test", "fqdn", "duba-nfws-sgwa01.corp.example.com"),
					resource.TestCheckResourceAttr("freeipa_host.test", "description", "sample description"),
					resource.TestCheckResourceAttr("freeipa_host.test", "id", "duba-nfws-sgwa01.corp.example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(fqdn string) string {
	return fmt.Sprintf(`
resource "freeipa_host" "test" {
  fqdn = %[1]q
}
`, fqdn)
}
