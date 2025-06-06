package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

// resourceSlackChannelUpdate handles updates to Slack channel resources.
func resourceSlackChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	api := meta.(*slack.Client)
	channelID := d.Id()

	// Handle name change
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		tflog.Info(ctx, fmt.Sprintf("Renaming Slack channel %s to '%s'", channelID, newName))

		_, err := api.RenameConversation(channelID, newName)
		if err != nil {
			return diag.Errorf("error renaming Slack channel: %s", err)
		}
	}

	// Handle privacy change (Slack does not support this)
	if d.HasChange("is_private") {
		oldRaw, newRaw := d.GetChange("is_private")
		oldPrivate := oldRaw.(bool)
		newPrivate := newRaw.(bool)

		if oldPrivate != newPrivate {
			return diag.Errorf(
				"Slack does not allow changing 'is_private' via API (Current: %v, Desired: %v). Please change manually in Slack.",
				oldPrivate, newPrivate,
			)
		}
	}

	// Handle member synchronization
	if d.HasChange("members") {
		_, newRaw := d.GetChange("members")
		newMembers := convertSchemaSetToStringSlice(newRaw.(*schema.Set))

		tflog.Debug(ctx, fmt.Sprintf("Syncing members for Slack channel %s", channelID))

		memberDiags := syncChannelMembers(api, channelID, newMembers)
		diags = append(diags, memberDiags...)
	}

	if d.HasChange("topic") {
		newTopic := d.Get("topic").(string)
		tflog.Debug(ctx, fmt.Sprintf("Updating topic for Slack channel %s to '%s'", channelID, newTopic))
		_, err := api.SetTopicOfConversation(channelID, newTopic)
		if err != nil {
			return diag.Errorf("error updating Slack channel topic: %s", err)
		}
	}

	if d.HasChange("purpose") {
		newPurpose := d.Get("purpose").(string)
		tflog.Debug(ctx, fmt.Sprintf("Updating purpose for Slack channel %s to '%s'", channelID, newPurpose))
		_, err := api.SetPurposeOfConversation(channelID, newPurpose)
		if err != nil {
			return diag.Errorf("error updating Slack channel purpose: %s", err)
		}
	}

	// Refresh state after update
	stateDiags := resourceSlackChannelRead(ctx, d, meta)
	diags = append(diags, stateDiags...)

	return diags
}

// convertInterfaceSliceToStringSlice safely converts []interface{} to []string
func convertInterfaceSliceToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = v.(string)
	}
	return result
}

func convertSchemaSetToStringSlice(set *schema.Set) []string {
	result := make([]string, 0, set.Len())
	for _, v := range set.List() {
		result = append(result, v.(string))
	}
	return result
}
