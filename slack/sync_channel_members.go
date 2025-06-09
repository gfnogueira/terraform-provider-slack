package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/slack-go/slack"
)

// syncChannelMembers ensures that all desired members are in the Slack channel.
// It adds missing members and warns about extra ones.
//
// Slack API does NOT support removing users via API, so extra users are only warned.
func syncChannelMembers(api *slack.Client, channelID string, desiredMembers []string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get current channel members from Slack
	currentMembers, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     1000,
	})
	if err != nil {
		return diag.Errorf("failed to retrieve current channel members from Slack: %s", err)
	}

	// Fetch bot ID (used to suppress its warning)
	auth, authErr := api.AuthTest()
	botID := ""
	if authErr == nil {
		botID = auth.UserID
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

	// Identify members to add
	var toAdd []string
	for _, user := range desiredMembers {
		if !currentSet[user] {
			toAdd = append(toAdd, user)
		}
	}
	if len(toAdd) > 0 {
		_, err := api.InviteUsersToConversation(channelID, toAdd...)
		if err != nil {
			return diag.Errorf("error adding members to Slack channel: %s", err)
		}
	}

	// Identify "extra" users, excluding the bot itself
	var extras []string
	for _, user := range currentMembers {
		if !desiredSet[user] && user != botID {
			// Try to resolve user info
			info, err := api.GetUserInfo(user)
			if err == nil {
				extras = append(extras, fmt.Sprintf("%s (%s)", info.Name, info.Profile.Email))
			} else {
				extras = append(extras, fmt.Sprintf("Unknown (%s)", user))
			}
		}
	}

	// Warn if extras exist (excluding bot)
	if len(extras) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Extra members found in Slack channel",
			Detail: fmt.Sprintf("Slack API does not support removing users via Terraform.\nExtra members present in the channel but not declared:\n- %s",
				strings.Join(extras, "\n- ")),
		})
	}

	return diags
}
