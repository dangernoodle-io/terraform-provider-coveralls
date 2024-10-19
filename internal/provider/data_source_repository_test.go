package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.coveralls_repository.test", "service", service),
					resource.TestCheckResourceAttr("data.coveralls_repository.test", "name", name),
				),
			},
		},
	})
}

var testAccExampleDataSourceConfig = fmt.Sprintf(`
data "coveralls_repository" "test" {
  service = "%s"
  name 	  = "%s"
}`, service, name)
