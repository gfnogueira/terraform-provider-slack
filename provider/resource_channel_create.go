package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannelCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	// Validate required members for private channels
	membersRaw := d.Get("members").([]interface{})
	var members []string
	for _, m := range membersRaw {
		members = append(members, m.(string))
	}
	if isPrivate && len(members) == 0 {
		return fmt.Errorf("private channels must have at least one member listed")
	}

	// Check if a channel with the same name already exists
	existingChannel, err := findChannelByName(api, name)
	if err != nil {
		return fmt.Errorf("error checking for existing channel: %w", err)
	}
	if existingChannel != nil {
		if existingChannel.IsArchived {
			log.Printf("[WARN] Channel '%s' exists but is archived. Reusing it — please unarchive it manually in Slack.", name)

			// ❗ wish: move to diag.Warning com CreateContext
		}
		d.SetId(existingChannel.ID)
		return resourceSlackChannelRead(d, meta)
	}

	// Create new channel
	params := slack.CreateConversationParams{
		ChannelName: name,
		IsPrivate:   isPrivate,
	}
	channel, err := api.CreateConversation(params)
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}
	d.SetId(channel.ID)

	// Add members to channel
	if len(members) > 0 {
		_, err := api.InviteUsersToConversation(channel.ID, members...)
		if err != nil {
			return fmt.Errorf("error adding members: %w", err)
		}
		log.Printf("[INFO] Members added to channel '%s': %v", name, members)
	}

	log.Printf("[INFO] Slack channel '%s' created successfully", name)
	return resourceSlackChannelRead(d, meta)
}
