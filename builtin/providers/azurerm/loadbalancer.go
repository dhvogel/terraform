package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/resource"
	"net/http"
)

func retrieveLoadbalancerById(loadBalancerId string, meta interface{}) (*network.LoadBalancer, bool, error) {
	loadBalancerClient := meta.(*ArmClient).loadBalancerClient

	id, err := parseAzureResourceID(loadBalancerId)
	if err != nil {
		return nil, false, err
	}
	name := id.Path["loadBalancers"]
	resGroup := id.ResourceGroup

	resp, err := loadBalancerClient.Get(resGroup, name, "")
	if err != nil {
		return nil, false, fmt.Errorf("Error making Read request on Azure Loadbalancer %s: %s", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}

	return &resp, true, nil
}

func mergeLoadbalancerConfig(oldLb *network.LoadBalancer) network.LoadBalancer {
	newLb := network.LoadBalancer{
		Name:     oldLb.Name,
		Etag:     oldLb.Etag,
		Location: oldLb.Location,
	}

	if oldLb.Tags != nil {
		newLb.Tags = oldLb.Tags
	}

	return newLb
}

func loadbalancerStateRefreshFunc(client *ArmClient, resourceGroupName string, loadbalancer string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.loadBalancerClient.Get(resourceGroupName, loadbalancer, "")
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in loadbalancerStateRefreshFunc to Azure ARM for Loadbalancer '%s' (RG: '%s'): %s", loadbalancer, resourceGroupName, err)
		}

		return res, *res.Properties.ProvisioningState, nil
	}
}
