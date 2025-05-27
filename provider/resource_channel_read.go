package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

// resourceSlackChannelRead fetches current Slack channel info from the API
// and updates the Terraform state accordingly.
func resourceSlackChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	channelID := d.Id()
	if channelID == "" {
		return diag.Errorf("no channel ID set")
	}

	// Request channel info from Slack
	info, err := api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		// Gracefully handle channel deletion (e.g. removed manually in Slack)
		if err.Error() == "channel_not_found" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Slack channel '%s' no longer exists", channelID),
				Detail:   "The channel was likely deleted manually in Slack. Terraform will now remove it from state.",
			})
			d.SetId("") // Mark resource for removal
			return diags
		}

		// Other errors should fail the plan
		return diag.Errorf("error reading channel info: %s", err)
	}

	// Set Terraform state values
	if err := d.Set("name", info.Name); err != nil {
		return diag.Errorf("error setting name: %s", err)
	}
	if err := d.Set("is_private", info.IsPrivate); err != nil {
		return diag.Errorf("error setting is_private: %s", err)
	}

	// Note: members cannot be reliably reloaded, as Slack API doesn't allow listing all members unless you're in the channel

	return diags
}
