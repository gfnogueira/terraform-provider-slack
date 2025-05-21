package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackUser() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			api := meta.(*slack.Client)
			email := d.Get("email").(string)

			user, err := api.GetUserByEmail(email)
			if err != nil {
				return fmt.Errorf("error searching for user by email: %w", err)
			}

			d.SetId(user.ID)
			d.Set("id", user.ID)
			d.Set("real_name", user.RealName)
			d.Set("email", email)
			return nil
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
