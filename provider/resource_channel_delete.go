package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)
	channelID := d.Id()

	err := api.ArchiveConversation(channelID)
	if err != nil {
		return fmt.Errorf("erro ao arquivar canal: %w", err)
	}

	d.SetId("")
	return nil
}
