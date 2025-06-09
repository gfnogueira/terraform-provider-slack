package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlackChannelCreate,
		ReadContext:   resourceSlackChannelRead,
		UpdateContext: resourceSlackChannelUpdate,
		DeleteContext: resourceSlackChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purpose (description) of the channel.",
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Topic shown at the top of the channel.",
			},
		},
	}
}
