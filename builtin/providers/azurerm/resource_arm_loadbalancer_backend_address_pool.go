package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/riviera/azure"
)

func resourceArmLoadbalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadbalancerBackendAddressPoolCreate,
		Read:   resourceArmLoadbalancerBackendAddressPoolRead,
		Delete: resourceArmLoadbalancerBackendAddressPoolDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backend_ip_configurations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"load_balancing_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceArmLoadbalancerBackendAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	lbClient := client.loadBalancerClient

	loadBalancer, exists, err := retrieveLoadbalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return err
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Loadbalancer %q not found. Refreshing from state", d.Get("name").(string))
		return nil
	}

	resGroup := d.Get("resource_group_name").(string)
	loadBalancerName := *loadBalancer.Name
	newLb := mergeLoadbalancerConfig(loadBalancer)
	newLb.Properties = &network.LoadBalancerPropertiesFormat{
		BackendAddressPools: expandAzureRmLoadbalancerBackendAddressPools(d),
	}

	_, err = lbClient.CreateOrUpdate(resGroup, loadBalancerName, newLb, make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := lbClient.Get(resGroup, loadBalancerName, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Loadbalancer %s (resource group %s) ID", loadBalancerName, resGroup)
	}

	d.SetId(*read.ID)

	log.Printf("[DEBUG] Waiting for LoadBalancer (%s) to become available", loadBalancerName)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: loadbalancerStateRefreshFunc(client, resGroup, loadBalancerName),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Loadbalancer (%s) to become available: %s", loadBalancerName, err)
	}

	return resourceArmLoadbalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadbalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	loadBalancer, exists, err := retrieveLoadbalancerById(d.Id(), meta)
	if err != nil {
		return err
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Loadbalancer %q not found. Refreshing from state", d.Get("name").(string))
		return nil
	}

	configs := *loadBalancer.Properties.BackendAddressPools
	for _, config := range configs {
		if *config.Name == d.Get("name").(string) {
			d.Set("name", config.Name)

			if config.Properties.BackendIPConfigurations != nil {
				backend_ip_configurations := make([]string, 0, len(*config.Properties.BackendIPConfigurations))
				for _, backendConfig := range *config.Properties.BackendIPConfigurations {
					backend_ip_configurations = append(backend_ip_configurations, *backendConfig.ID)
				}

				d.Set("backend_ip_configurations", backend_ip_configurations)
			}

			if config.Properties.LoadBalancingRules != nil {
				load_balancing_rules := make([]string, 0, len(*config.Properties.LoadBalancingRules))
				for _, rule := range *config.Properties.LoadBalancingRules {
					load_balancing_rules = append(load_balancing_rules, *rule.ID)
				}

				d.Set("backend_ip_configurations", load_balancing_rules)
			}

			break
		}
	}

	return nil
}

func resourceArmLoadbalancerBackendAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func expandAzureRmLoadbalancerBackendAddressPools(d *schema.ResourceData) *[]network.BackendAddressPool {
	backendPools := make([]network.BackendAddressPool, 0)

	backendPool := network.BackendAddressPool{
		Name: azure.String(d.Get("name").(string)),
	}

	backendPools = append(backendPools, backendPool)

	return &backendPools
}
