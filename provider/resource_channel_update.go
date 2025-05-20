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
		newPrivate := d.Get("is_private").(bool)
		return fmt.Errorf("⚠️ Slack API does not allow changing `is_private` via code (desired: %v)", newPrivate)
	}

	return resourceSlackChannelRead(d, meta)
}
