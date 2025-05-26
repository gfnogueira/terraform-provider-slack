package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)
	channelID := d.Id()

	// Attempt to join the channel before archiving
	_, _, _, err := api.JoinConversation(channelID)
	if err != nil {
		log.Printf("[WARN] Unable to join channel '%s': %v. Proceeding to attempt archive anyway.", channelID, err)
	}

	// Attempt to archive the channel
	err = api.ArchiveConversation(channelID)
	if err != nil {
		slackErr := err.Error()

		switch slackErr {
		case "not_in_channel":
			log.Printf("[WARN] Cannot archive channel '%s' because the bot is not a member. Please archive it manually if needed.", channelID)
			return nil
		case "channel_not_found":
			log.Printf("[WARN] Channel '%s' was not found. Assuming it was manually deleted or the ID is invalid.", channelID)
			return nil
		default:
			return fmt.Errorf("error archiving channel '%s': %w", channelID, err)
		}
	}

	log.Printf("[INFO] Slack channel '%s' archived successfully", channelID)
	return nil
}
