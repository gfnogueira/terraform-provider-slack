package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceSlackUsersGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlackUsersGroupRead,

		Schema: map[string]*schema.Schema{
			"emails": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceSlackUsersGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	var diags diag.Diagnostics

	rawEmails := d.Get("emails").([]interface{})
	var userIDs []string
	var missing []string

	for _, raw := range rawEmails {
		email := raw.(string)
		user, err := api.GetUserByEmail(email)
		if err != nil {
			missing = append(missing, email)
			continue
		}
		userIDs = append(userIDs, user.ID)
	}

	if len(missing) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Some users not found in Slack",
			Detail:   fmt.Sprintf("The following emails could not be resolved: %v", missing),
		})
	}

	sort.Strings(userIDs)
	d.Set("ids", schema.NewSet(schema.HashString, convertStringSliceToInterface(userIDs)))
	d.SetId(fmt.Sprintf("group-%x", hashStringSlice(userIDs)))

	return diags
}

func hashStringSlice(slice []string) string {
	h := sha1.New()
	for _, s := range slice {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func convertStringSliceToInterface(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
