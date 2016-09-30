package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmLoadbalancerNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadbalancerNatRuleCreate,
		Read:   resourceArmLoadbalancerNatRuleRead,
		Update: resourceArmLoadbalancerNatRuleCreate,
		Delete: resourceArmLoadbalancerNatRuleDelete,

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

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"backend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadbalancerNatRuleCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmLoadbalancerNatRuleRead(d, meta)
}

func resourceArmLoadbalancerNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmLoadbalancerNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
