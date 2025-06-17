package slack

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

// findChannelByName searches for a Slack channel by its name.
// Returns the channel object if found (including archived), or nil if not found.
// Logs and diagnostics are returned as formatted strings for flexibility.
func findChannelByName(ctx context.Context, api *slack.Client, name string) (*slack.Channel, []string, error) {
	var debugLogs []string
	cursor := ""

	for {
		debugLogs = append(debugLogs, fmt.Sprintf("[DEBUG] Querying Slack API for channels (cursor: '%s')", cursor))

		channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
			ExcludeArchived: false,
			Limit:           1000,
			Cursor:          cursor,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return nil, debugLogs, fmt.Errorf("failed to list Slack channels: %w", err)
		}

		debugLogs = append(debugLogs, fmt.Sprintf("[DEBUG] Slack API returned %d channels", len(channels)))

		for _, c := range channels {
			debugLogs = append(debugLogs, fmt.Sprintf("[DEBUG] Found channel: %s (ID: %s, Archived: %v)", c.Name, c.ID, c.IsArchived))
			if c.Name == name {
				debugLogs = append(debugLogs, fmt.Sprintf("[INFO] Matched requested channel '%s' (Archived: %v)", name, c.IsArchived))
				return &c, debugLogs, nil
			}
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	debugLogs = append(debugLogs, fmt.Sprintf("[INFO] Channel '%s' not found in Slack API (even including archived)", name))
	return nil, debugLogs, nil
}
