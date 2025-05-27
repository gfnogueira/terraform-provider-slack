package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	api := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	// Validate required members for private channels
	membersRaw := d.Get("members").([]interface{})
	var members []string
	for _, m := range membersRaw {
		members = append(members, m.(string))
	}
	if isPrivate && len(members) == 0 {
		return diag.Errorf("private channels must have at least one member listed")
	}

	// Check if a channel with the same name already exists
	existingChannel, err := findChannelByName(api, name)
	if err != nil {
		return diag.Errorf("error checking for existing channel: %s", err)
	}

	if existingChannel != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Channel '%s' already exists", name),
			Detail:   fmt.Sprintf("Terraform will reuse this channel. Archived status: %v. Please unarchive it in Slack if necessary.", existingChannel.IsArchived),
		})

		d.SetId(existingChannel.ID)
		diags = append(diags, resourceSlackChannelRead(ctx, d, meta)...)
		return diags
	}

	// Create new channel
	params := slack.CreateConversationParams{
		ChannelName: name,
		IsPrivate:   isPrivate,
	}
	channel, err := api.CreateConversation(params)
	if err != nil {
		return diag.Errorf("error creating channel: %s", err)
	}
	d.SetId(channel.ID)

	// Add members
	if len(members) > 0 {
		_, err := api.InviteUsersToConversation(channel.ID, members...)
		if err != nil {
			return diag.Errorf("error adding members: %s", err)
		}
	}

	diags = append(diags, resourceSlackChannelRead(ctx, d, meta)...)
	return diags
}
