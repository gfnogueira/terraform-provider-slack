package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

// resourceSlackChannelRead fetches current Slack channel info from the Slack API
// and updates the Terraform state. It gracefully handles deleted channels and
// filters out the bot user from the members list to prevent unwanted diffs.
func resourceSlackChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	channelID := d.Id()
	if channelID == "" {
		return diag.Errorf("no channel ID set")
	}

	// Fetch basic channel info
	info, err := api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		// Handle channel manually deleted in Slack
		if err.Error() == "channel_not_found" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Slack channel '%s' no longer exists", channelID),
				Detail:   "The channel was likely deleted manually in Slack. Terraform will now remove it from state.",
			})
			d.SetId("") // Remove from Terraform state
			return diags
		}
		return diag.Errorf("error reading channel info: %s", err)
	}

	// Set core attributes
	if err := d.Set("name", info.Name); err != nil {
		return diag.Errorf("error setting 'name': %s", err)
	}
	if err := d.Set("is_private", info.IsPrivate); err != nil {
		return diag.Errorf("error setting 'is_private': %s", err)
	}
	if err := d.Set("purpose", info.Purpose.Value); err != nil {
		return diag.Errorf("error setting 'purpose': %s", err)
	}
	
	if err := d.Set("topic", info.Topic.Value); err != nil {
		return diag.Errorf("error setting 'topic': %s", err)
	}

	// Fetch channel members
	members, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     1000,
	})
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("failed to fetch members for channel '%s': %s", channelID, err))
		// Proceed without setting "members"
		return diags
	}

	// Fetch bot's own user ID to exclude from member list
	authResp, err := api.AuthTest()
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("failed to identify bot user ID: %s", err))
	} else {
		botUserID := authResp.UserID
		filtered := make([]string, 0, len(members))
		for _, m := range members {
			if m != botUserID {
				filtered = append(filtered, m)
			}
		}
		members = filtered
	}

	if err := d.Set("members", members); err != nil {
		return diag.Errorf("error setting 'members': %s", err)
	}

	return diags
}
