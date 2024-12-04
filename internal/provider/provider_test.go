package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
//var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//	"freeipa": providerserver.NewProtocol6WithError(New("test")()),
//}

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"freeipa": func() (tfprotov5.ProviderServer, error) {
		return tf6to5server.DowngradeServer(
			context.Background(),
			providerserver.NewProtocol6(New("test")()),
		)
	},
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
