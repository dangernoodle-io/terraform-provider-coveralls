package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-coveralls/internal/provider"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/dangernoode-io/coveralls",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
