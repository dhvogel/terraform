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

func resourceArmLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadbalancerCreate,
		Read:   resourceArmLoadbalancerRead,
		Update: resourceArmLoadbalancerCreate,
		Delete: resourceArmLoadbalancerDelete,

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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	loadBalancerClient := client.loadBalancerClient

	log.Printf("[INFO] preparing arguments for Azure ARM Loadbalancer creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	loadbalancer := network.LoadBalancer{
		Name:     azure.String(name),
		Location: azure.String(location),
		Tags:     expandedTags,
	}

	_, err := loadBalancerClient.CreateOrUpdate(resGroup, name, loadbalancer, make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := loadBalancerClient.Get(resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Loadbalancer %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	log.Printf("[DEBUG] Waiting for LoadBalancer (%s) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: loadbalancerStateRefreshFunc(client, resGroup, name),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Loadbalancer (%s) to become available: %s", name, err)
	}

	return resourceArmLoadbalancerRead(d, meta)
}

func resourceArmLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	loadBalancer, exists, err := retrieveLoadbalancerById(d.Id(), meta)
	if err != nil {
		return err
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Loadbalancer %q not found. Refreshing from state", d.Get("name").(string))
		return nil
	}

	flattenAndSetTags(d, loadBalancer.Tags)

	return nil
}

func resourceArmLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	loadBalancerClient := meta.(*ArmClient).loadBalancerClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["loadBalancers"]

	_, err = loadBalancerClient.Delete(resGroup, name, make(chan struct{}))

	return err
}
