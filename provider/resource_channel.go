package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceSlackChannelCreate,
		Read:   resourceSlackChannelRead,
		Update: resourceSlackChannelUpdate,
		Delete: resourceSlackChannelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"members": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
