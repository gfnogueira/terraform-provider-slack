package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	channelID := d.Id()

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		_, err := api.RenameConversation(channelID, newName)
		if err != nil {
			return diag.Errorf("error renaming channel: %s", err)
		}
	}

	if d.HasChange("is_private") {
		old, new := d.GetChange("is_private")
		oldPrivate := old.(bool)
		newPrivate := new.(bool)

		if oldPrivate != newPrivate {
			return diag.Errorf(
				"Slack does not allow changing 'is_private' via API (Current: %v, desired: %v). Please change manually in Slack.",
				oldPrivate, newPrivate,
			)
		}
	}

	return resourceSlackChannelRead(ctx, d, meta)
}
