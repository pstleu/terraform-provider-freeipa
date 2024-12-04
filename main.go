package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-freeipa/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	//opts := providerserver.ServeOpts{
	//	// TODO: Update this string with the published name of your provider.
	//	// Also update the tfplugindocs generate command to either remove the
	//	// -provider-name flag or set its value to the updated provider name.
	//	Address: "registry.terraform.io/mashanm/freeipa",
	//	Debug:   debug,
	//}

	downgradedFrameworkProvider, err := tf6to5server.DowngradeServer(
		context.Background(),
		providerserver.NewProtocol6(provider.New(version)()),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = tf5server.Serve(
		"registry.terraform.io/mashanm/freeipa",
		func() tfprotov5.ProviderServer {
			return downgradedFrameworkProvider
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}
}
