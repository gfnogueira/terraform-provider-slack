package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/slack-go/slack"
)

// syncChannelMembers ensures that all desired members are in the Slack channel.
// It adds missing members and warns about extra ones.
//
// Slack API does NOT support removing users from channels via API, so we emit a warning for extra members.
//
// This function also detects the bot's own user ID (via AuthTest) and excludes it from extra-member warnings.
func syncChannelMembers(api *slack.Client, channelID string, desiredMembers []string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get current members from Slack
	currentMembers, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     1000,
	})
	if err != nil {
		return diag.Errorf("failed to retrieve current members from Slack: %s", err)
	}

	// Get bot's own user ID to exclude it from extra warnings
	authInfo, err := api.AuthTest()
	botID := ""
	if err == nil {
		botID = authInfo.UserID
	}

	// Create lookup sets
	currentSet := make(map[string]bool)
	for _, user := range currentMembers {
		currentSet[user] = true
	}

	desiredSet := make(map[string]bool)
	for _, user := range desiredMembers {
		desiredSet[user] = true
	}

	// Determine users to invite
	var toAdd []string
	for _, user := range desiredMembers {
		if !currentSet[user] {
			toAdd = append(toAdd, user)
		}
	}

	// Invite missing users
	if len(toAdd) > 0 {
		_, err := api.InviteUsersToConversation(channelID, toAdd...)
		if err != nil {
			return diag.Errorf("error adding members to channel: %s", err)
		}
	}

	// Identify "extra" users, excluding the bot itself
	var extras []string
	for _, user := range currentMembers {
		if user != botID && !desiredSet[user] {
			extras = append(extras, user)
		}
	}

	// Emit warning for extras
	if len(extras) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Extra members found in Slack channel",
			Detail:   fmt.Sprintf("Slack API does not support removing users via Terraform. Extra members: %v", extras),
		})
	}

	return diags
}
