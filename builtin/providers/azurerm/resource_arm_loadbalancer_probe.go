package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmLoadbalancerProbe() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadbalancerProbeCreate,
		Read:   resourceArmLoadbalancerProbeRead,
		Update: resourceArmLoadbalancerProbeCreate,
		Delete: resourceArmLoadbalancerProbeDelete,

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

			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"request_path": {
				Type:     schema.TypeString,
				Required: true,
			},

			"interval_in_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"number_of_probes": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"load_balance_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceArmLoadbalancerProbeCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmLoadbalancerProbeRead(d, meta)
}

func resourceArmLoadbalancerProbeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmLoadbalancerProbeDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
