package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	service = "github"
	name    = "jgangemi/terraform-coveralls-test"

	providerConfig = `
provider "coveralls" {
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"coveralls": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("COVERALLS_API_TOKEN") == "" {
		t.Fatal("COVERALLS_API_TOKEN must be set for acceptance tests")
	}
}
