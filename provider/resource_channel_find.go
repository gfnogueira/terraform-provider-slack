package provider

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

// findChannelByName searches for a Slack channel by its name.
// It returns the channel object if found (even if archived), or nil if not found.
// Returns an error if the Slack API fails.
func findChannelByName(api *slack.Client, name string) (*slack.Channel, error) {
	cursor := ""

	for {
		log.Printf("[DEBUG] Querying Slack API for channels (cursor: '%s')", cursor)

		channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
			ExcludeArchived: false, // VERY IMPORTANT: includes archived channels
			Limit:           1000,
			Cursor:          cursor,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list Slack channels: %w", err)
		}

		log.Printf("[DEBUG] Slack API returned %d channels", len(channels))

		for _, c := range channels {
			log.Printf("[DEBUG] Found channel: %s (ID: %s, Archived: %v)", c.Name, c.ID, c.IsArchived)

			if c.Name == name {
				log.Printf("[INFO] Matched requested channel '%s' (Archived: %v)", name, c.IsArchived)
				return &c, nil
			}
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	log.Printf("[INFO] Channel '%s' not found in Slack API (even including archived)", name)
	return nil, nil
}
