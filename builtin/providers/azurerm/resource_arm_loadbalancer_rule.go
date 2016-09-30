package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmLoadbalancerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadbalancerRuleCreate,
		Read:   resourceArmLoadbalancerRuleRead,
		Update: resourceArmLoadbalancerRuleCreate,
		Delete: resourceArmLoadbalancerRuleDelete,

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

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"backend_address_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},

			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"backend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"probe_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"load_distribution": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadbalancerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmLoadbalancerRuleRead(d, meta)
}

func resourceArmLoadbalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmLoadbalancerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
