package provider

import (
	"fmt"

	"github.com/slack-go/slack"
)

// findChannelByName searches for a Slack channel by its name.
// It returns the channel object if found (even if archived), or nil if not found.
// Returns an error if the Slack API fails.
func findChannelByName(api *slack.Client, name string) (*slack.Channel, error) {
	cursor := ""

	for {
		channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
			ExcludeArchived: false,
			Limit:           1000,
			Cursor:          cursor,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list Slack channels: %w", err)
		}

		for _, c := range channels {
			if c.Name == name {
				return &c, nil
			}
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return nil, nil
}
