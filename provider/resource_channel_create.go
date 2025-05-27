package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	// Validate members list
	membersRaw := d.Get("members").([]interface{})
	var members []string
	for _, m := range membersRaw {
		members = append(members, m.(string))
	}

	if isPrivate && len(members) == 0 {
		return diag.Errorf("private channels must have at least one member listed")
	}
	// Check if a channel with the same name already exists
	existingChannel, debugLogs, err := findChannelByName(ctx, api, name)
	for _, log := range debugLogs {
		tflog.Debug(ctx, log)
	}
	if err != nil {
		return diag.Errorf("error checking for existing channel: %s", err)
	}

	if existingChannel != nil {
		diagMsg := fmt.Sprintf("Channel '%s' already exists", name)
		detail := fmt.Sprintf("Terraform will reuse this channel. Archived status: %v. Please unarchive it in Slack if necessary.", existingChannel.IsArchived)

		tflog.Warn(ctx, fmt.Sprintf("%s — %s", diagMsg, detail))
		d.SetId(existingChannel.ID)

		// Return warning via diag (visible in CLI)
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  diagMsg,
			Detail:   detail,
		}}
	}

	// Create channel
	params := slack.CreateConversationParams{
		ChannelName: name,
		IsPrivate:   isPrivate,
	}
	channel, err := api.CreateConversation(params)
	if err != nil {
		return diag.Errorf("error creating channel: %s", err)
	}

	d.SetId(channel.ID)

	// Add members if any
	if len(members) > 0 {
		_, err := api.InviteUsersToConversation(channel.ID, members...)
		if err != nil {
			return diag.Errorf("error adding members: %s", err)
		}
		tflog.Info(ctx, fmt.Sprintf("Members added to channel '%s': %v", name, members))
	}

	tflog.Info(ctx, fmt.Sprintf("Slack channel '%s' created successfully", name))
	return resourceSlackChannelRead(ctx, d, meta)
}
