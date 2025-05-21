package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)
	channelID := d.Id()

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		_, err := api.RenameConversation(channelID, newName)
		if err != nil {
			return fmt.Errorf("error renaming channel: %w", err)
		}
	}

	if d.HasChange("is_private") {
		oldRaw, newRaw := d.GetChange("is_private")
		oldPrivate := oldRaw.(bool)
		newPrivate := newRaw.(bool)

		return fmt.Errorf(
			"⚠️ Changing 'is_private' from %v to %v is not supported by the Slack API. Please update the channel visibility manually in Slack before applying this change via Terraform",
			oldPrivate, newPrivate,
		)
	}

	return resourceSlackChannelRead(d, meta)
}
