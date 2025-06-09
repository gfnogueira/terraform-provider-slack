package provider

import (
	"fmt"

	"github.com/slack-go/slack"
)

// syncChannelMembers ensures the desired list of members are added to the channel.
func syncChannelMembers(api *slack.Client, channelID string, current []string, desired []string) error {
	var toAdd []string

	// Detect members to add
	existing := make(map[string]bool)
	for _, m := range current {
		existing[m] = true
	}
	for _, m := range desired {
		if !existing[m] {
			toAdd = append(toAdd, m)
		}
	}

	if len(toAdd) > 0 {
		_, err := api.InviteUsersToConversation(channelID, toAdd...)
		if err != nil {
			return fmt.Errorf("error inviting users to channel: %w", err)
		}
	}

	return nil
}
