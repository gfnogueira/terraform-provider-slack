package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)
	channelID := d.Id()

	input := &slack.GetConversationInfoInput{
		ChannelID: channelID,
	}

	info, err := api.GetConversationInfo(input)
	if err != nil {
		return fmt.Errorf("error reading channel: %w", err)
	}

	d.Set("name", info.Name)
	d.Set("is_private", info.IsPrivate)

	return nil
}
