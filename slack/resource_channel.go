package slack

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the channel (without #).",
			},
			"is_private": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the channel is private (true) or public (false).",
			},
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of user IDs to add to the channel.",
			},
			"strict_members": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, Terraform will detect drift when users are manually added to the channel. If false (default), Terraform only manages the declared members and ignores manually added users.",
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
