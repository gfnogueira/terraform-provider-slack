package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackChannelsRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter channels whose names start with this prefix.",
			},
			"include_archived": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Include archived channels in the result.",
			},
			"is_private": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set, filters by channel privacy (true = private, false = public).",
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_archived": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSlackChannelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	var allChannels []map[string]interface{}
	prefix := d.Get("prefix").(string)
	includeArchived := d.Get("include_archived").(bool)
	filterPrivate, hasPrivacy := d.GetOk("is_private")

	cursor := ""
	for {
		channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
			ExcludeArchived: !includeArchived,
			Limit:           1000,
			Cursor:          cursor,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return diag.Errorf("failed to list Slack channels: %s", err)
		}

		for _, c := range channels {
			// Filter by prefix
			if prefix != "" && !strings.HasPrefix(c.Name, prefix) {
				continue
			}
			// Filter by privacy
			if hasPrivacy && c.IsPrivate != filterPrivate.(bool) {
				continue
			}
			allChannels = append(allChannels, map[string]interface{}{
				"id":          c.ID,
				"name":        c.Name,
				"is_private":  c.IsPrivate,
				"is_archived": c.IsArchived,
			})
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	d.Set("channels", allChannels)
	d.SetId(fmt.Sprintf("channels-%d", len(allChannels)))

	return diags
}
