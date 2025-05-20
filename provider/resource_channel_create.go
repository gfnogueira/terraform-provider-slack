package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

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
		return fmt.Errorf("erro ao criar canal: %w", err)
	}

	d.SetId(channel.ID)

	// Adiciona membros, se houver
	membersRaw := d.Get("members").([]interface{})
	var members []string
	for _, m := range membersRaw {
		members = append(members, m.(string))
	}

	if isPrivate && len(members) == 0 {
		return fmt.Errorf("canais privados precisam ter ao menos um membro listado")
	}

	if len(members) > 0 {
		_, err := api.InviteUsersToConversation(channel.ID, members...)
		if err != nil {
			return fmt.Errorf("erro ao adicionar membros: %w", err)
		}
	}

	return resourceSlackChannelRead(d, meta)
}
