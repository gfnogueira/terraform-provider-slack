package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackUserRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_owner": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_bot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSlackUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	email := d.Get("email").(string)
	tflog.Info(ctx, fmt.Sprintf("Searching for Slack user with email: %s", email))

	user, err := api.GetUserByEmail(email)
	if err != nil {
		return diag.Errorf("error retrieving Slack user by email '%s': %s", email, err)
	}

	d.SetId(user.ID)

	set := func(key string, val interface{}) {
		if err := d.Set(key, val); err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to set '%s': %s", key, err))
		}
	}

	set("id", user.ID)
	set("real_name", user.RealName)
	set("display_name", user.Profile.DisplayName)
	set("email", email)
	set("is_admin", user.IsAdmin)
	set("is_owner", user.IsOwner)
	set("is_bot", user.IsBot)
	set("timezone", user.TZ)

	return diags
}