package slack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SLACK_TOKEN", nil),
				Description: "Slack API Token (starts with xoxb-)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"slack_channel":   resourceSlackChannel(),
			"slack_usergroup": resourceSlackUsergroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"slack_user":        dataSourceSlackUser(),
			"slack_users":       dataSourceSlackUsers(),
			"slack_users_group": dataSourceSlackUsersGroup(),
			"slack_channel":     dataSourceSlackChannel(),
			"slack_channels":    dataSourceSlackChannels(),
		},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			token := d.Get("token").(string)
			if token == "" {
				return nil, diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Slack token is missing",
						Detail:   "Please provide a Slack token using the 'token' argument or the SLACK_TOKEN environment variable.",
					},
				}
			}

			client := slack.New(token)

			//Validate token authenticity
			_, err := client.AuthTest()
			if err != nil {
				return nil, diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid Slack token",
						Detail:   "Authentication with Slack API failed. Please check if your token is valid and has the correct permissions.",
					},
				}
			}

			return client, nil
		},
	}
}
