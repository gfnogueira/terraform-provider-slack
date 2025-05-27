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
	var diags diag.Diagnostics
	api := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	// Parse members
	membersRaw := d.Get("members").([]interface{})
	var members []string
	for _, m := range membersRaw {
		members = append(members, m.(string))
	}

	// Validate members for private channels
	if isPrivate && len(members) == 0 {
		return diag.Errorf("private channels must have at least one member listed")
	}

	// Check if channel already exists (including archived)
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

		// Sync members (automatically checks current vs desired)
		memberDiags := syncChannelMembers(api, existingChannel.ID, members)
		diags = append(diags, memberDiags...)

		// Return warning to user
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  diagMsg,
			Detail:   detail,
		})
		return diags
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

	// Sync members
	memberDiags := syncChannelMembers(api, channel.ID, members)
	diags = append(diags, memberDiags...)

	tflog.Info(ctx, fmt.Sprintf("Slack channel '%s' created successfully", name))
	return diags
}
