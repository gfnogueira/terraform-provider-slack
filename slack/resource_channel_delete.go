package slack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	api := meta.(*slack.Client)
	channelID := d.Id()

	// Attempt to join the channel (required before archiving)
	_, _, _, err := api.JoinConversation(channelID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unable to join Slack channel",
			Detail:   fmt.Sprintf("Bot could not join channel '%s': %v. Terraform will attempt to archive it anyway.", channelID, err),
		})
	}

	// Attempt to archive the channel
	err = api.ArchiveConversation(channelID)
	if err != nil {
		switch err.Error() {
		case "not_in_channel":
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bot is not a member of the Slack channel",
				Detail:   fmt.Sprintf("Terraform could not archive the channel '%s' because the bot is not a member. Please archive it manually in Slack.", channelID),
			})
			return diags
		case "channel_not_found":
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Slack channel not found",
				Detail:   fmt.Sprintf("Terraform could not find the channel '%s'. It may have been deleted manually.", channelID),
			})
			return diags
		default:
			return diag.FromErr(fmt.Errorf("error archiving Slack channel '%s': %w", channelID, err))
		}
	}

	// No need to log success — Terraform CLI handles positive feedback itself
	return diags
}
