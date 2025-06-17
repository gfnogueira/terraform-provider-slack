package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackUsersRead,
		Schema: map[string]*schema.Schema{
			"domain_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If provided, only users whose email ends with this domain will be included (e.g. '@empresa.com')",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            {Type: schema.TypeString, Computed: true},
						"real_name":     {Type: schema.TypeString, Computed: true},
						"display_name":  {Type: schema.TypeString, Computed: true},
						"email":         {Type: schema.TypeString, Computed: true},
						"is_admin":      {Type: schema.TypeBool, Computed: true},
						"is_owner":      {Type: schema.TypeBool, Computed: true},
						"is_bot":        {Type: schema.TypeBool, Computed: true},
						"timezone":      {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceSlackUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	domainFilter := d.Get("domain_filter").(string)
	if domainFilter != "" {
		tflog.Info(ctx, fmt.Sprintf("Applying domain filter: '%s'", domainFilter))
	} else {
		tflog.Info(ctx, "No domain filter applied, returning all users")
	}

	users, err := api.GetUsers()
	if err != nil {
		return diag.Errorf("error retrieving Slack users: %s", err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Fetched %d users from Slack API", len(users)))

	var filtered []map[string]interface{}
	for _, user := range users {
		if user.Deleted || user.IsBot {
			continue
		}
		email := user.Profile.Email
		if domainFilter == "" || (email != "" && strings.HasSuffix(email, domainFilter)) {
			filtered = append(filtered, map[string]interface{}{
				"id":           user.ID,
				"real_name":    user.RealName,
				"display_name": user.Profile.DisplayName,
				"email":        email,
				"is_admin":     user.IsAdmin,
				"is_owner":     user.IsOwner,
				"is_bot":       user.IsBot,
				"timezone":     user.TZ,
			})
		}
	}
	tflog.Info(ctx, fmt.Sprintf("Returning %d users after applying filter", len(filtered)))

	if err := d.Set("users", filtered); err != nil {
		return diag.Errorf("error setting users list: %s", err)
	}

	d.SetId(fmt.Sprintf("slack-users-%d", len(filtered)))
	return diags
}