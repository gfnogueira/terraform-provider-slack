package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackChannelRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Slack channel to look up.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Slack channel.",
			},
			"is_private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the channel is private.",
			},
		},
	}
}

func dataSourceSlackChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	name := d.Get("name").(string)

	channel, debugLogs, err := findChannelByName(ctx, api, name)
	for _, msg := range debugLogs {
		// Optional: debug log
		tflog.Debug(ctx, msg)
	}

	if err != nil {
		return diag.Errorf("error searching for Slack channel: %s", err)
	}

	if channel == nil {
		return diag.Errorf("Slack channel '%s' not found", name)
	}

	d.SetId(channel.ID)
	_ = d.Set("is_private", channel.IsPrivate)

	return nil
}
