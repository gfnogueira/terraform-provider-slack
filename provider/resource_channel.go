package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceSlackChannelCreate,
		Read:   resourceSlackChannelRead,
		Update: resourceSlackChannelUpdate,
		Delete: resourceSlackChannelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSlackChannelCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	params := slack.CreateConversationParams{
		ChannelName: name,
		IsPrivate:   isPrivate,
	}

	channel, err := api.CreateConversation(params)
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}

	d.SetId(channel.ID)
	return resourceSlackChannelRead(d, meta)
}

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

func resourceSlackChannelDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*slack.Client)
	channelID := d.Id()

	err := api.ArchiveConversation(channelID)
	if err != nil {
		return fmt.Errorf("error archiving channel: %w", err)
	}

	d.SetId("")
	return nil
}
