package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMLoadbalancerFrontendIpConfig_basic(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadbalancerFrontendIpConfig_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadbalancerExists("azurerm_lb.test"),
					//resource.TestCheckResourceAttr(
					//	"azurerm_lb_frontend_ip_config.test", "name", ri),
				),
			},
		},
	})
}

func testAccAzureRMLoadbalancerFrontendIpConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestrg-%d"
    location = "West US"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    tags {
    	Environment = "production"
    	Purpose = "AcceptanceTests"
    }
}

resource "azurerm_lb_frontend_ip_config" "test" {
    name = "%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    loadbalancer_id = "${azurerm_lb.test.id}"
}`, rInt, rInt, rInt)
}
