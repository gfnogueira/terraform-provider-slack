package slack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func resourceSlackUsergroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlackUsergroupCreate,
		ReadContext:   resourceSlackUsergroupRead,
		UpdateContext: resourceSlackUsergroupUpdate,
		DeleteContext: resourceSlackUsergroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"handle": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The handle/mention name for the usergroup (e.g., 'developers' for @developers)",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The display name for the usergroup",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the usergroup",
			},
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of user IDs that are members of this usergroup",
			},
			// Computed fields
			"team_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The team ID this usergroup belongs to",
			},
		},
	}
}

func resourceSlackUsergroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)

	handle := d.Get("handle").(string)
	name := d.Get("name").(string)
	if name == "" {
		name = handle // Use handle as name if not provided
	}
	description := d.Get("description").(string)

	tflog.Debug(ctx, fmt.Sprintf("Creating usergroup with handle: %s", handle))

	// Create the usergroup
	userGroup := slack.UserGroup{
		Handle:      handle,
		Name:        name,
		Description: description,
	}

	createdGroup, err := api.CreateUserGroupContext(ctx, userGroup)
	if err != nil {
		return diag.Errorf("error creating usergroup: %s", err)
	}

	d.SetId(createdGroup.ID)
	tflog.Info(ctx, fmt.Sprintf("Usergroup created with ID: %s", createdGroup.ID))

	// Update members if provided
	if membersRaw, ok := d.GetOk("members"); ok {
		members := convertSchemaSetToStringSlice(membersRaw.(*schema.Set))
		if len(members) > 0 {
			memberDiags := updateUsergroupMembers(ctx, api, createdGroup.ID, members)
			if memberDiags.HasError() {
				return memberDiags
			}
		}
	}

	return resourceSlackUsergroupRead(ctx, d, meta)
}

func resourceSlackUsergroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	usergroupID := d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Reading usergroup: %s", usergroupID))

	// Get all usergroups and find ours
	usergroups, err := api.GetUserGroupsContext(ctx,
		slack.GetUserGroupsOptionIncludeUsers(true),
		slack.GetUserGroupsOptionIncludeDisabled(true),
	)
	if err != nil {
		return diag.Errorf("error reading usergroups: %s", err)
	}

	var usergroup *slack.UserGroup
	for i := range usergroups {
		if usergroups[i].ID == usergroupID {
			usergroup = &usergroups[i]
			break
		}
	}

	if usergroup == nil {
		tflog.Warn(ctx, fmt.Sprintf("Usergroup %s not found, removing from state", usergroupID))
		d.SetId("")
		return nil
	}

	// Set all the fields
	if err := d.Set("handle", usergroup.Handle); err != nil {
		return diag.Errorf("error setting handle: %s", err)
	}
	if err := d.Set("name", usergroup.Name); err != nil {
		return diag.Errorf("error setting name: %s", err)
	}
	if err := d.Set("description", usergroup.Description); err != nil {
		return diag.Errorf("error setting description: %s", err)
	}
	if err := d.Set("team_id", usergroup.TeamID); err != nil {
		return diag.Errorf("error setting team_id: %s", err)
	}

	// Set members
	if len(usergroup.Users) > 0 {
		membersSet := schema.NewSet(schema.HashString, convertStringSliceToInterface(usergroup.Users))
		if err := d.Set("members", membersSet); err != nil {
			return diag.Errorf("error setting members: %s", err)
		}
	}

	return nil
}

func resourceSlackUsergroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	usergroupID := d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Updating usergroup: %s", usergroupID))

	// Update name, handle, or description if changed
	if d.HasChanges("name", "handle", "description") {
		var options []slack.UpdateUserGroupsOption

		if d.HasChange("name") {
			name := d.Get("name").(string)
			options = append(options, slack.UpdateUserGroupsOptionName(name))
		}
		if d.HasChange("handle") {
			handle := d.Get("handle").(string)
			options = append(options, slack.UpdateUserGroupsOptionHandle(handle))
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			options = append(options, slack.UpdateUserGroupsOptionDescription(&description))
		}

		_, err := api.UpdateUserGroupContext(ctx, usergroupID, options...)
		if err != nil {
			return diag.Errorf("error updating usergroup: %s", err)
		}
		tflog.Info(ctx, "Usergroup metadata updated successfully")
	}

	// Update members if changed
	if d.HasChange("members") {
		var members []string
		if membersRaw, ok := d.GetOk("members"); ok {
			members = convertSchemaSetToStringSlice(membersRaw.(*schema.Set))
		}

		memberDiags := updateUsergroupMembers(ctx, api, usergroupID, members)
		if memberDiags.HasError() {
			return memberDiags
		}
	}

	return resourceSlackUsergroupRead(ctx, d, meta)
}

func resourceSlackUsergroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*slack.Client)
	usergroupID := d.Id()

	tflog.Info(ctx, fmt.Sprintf("Disabling usergroup: %s", usergroupID))

	// Slack doesn't allow deleting usergroups, only disabling them
	_, err := api.DisableUserGroupContext(ctx, usergroupID)
	if err != nil {
		return diag.Errorf("error disabling usergroup: %s", err)
	}

	d.SetId("")
	tflog.Info(ctx, "Usergroup disabled successfully")
	return nil
}

// Helper function to update usergroup members
func updateUsergroupMembers(ctx context.Context, api *slack.Client, usergroupID string, members []string) diag.Diagnostics {
	tflog.Debug(ctx, fmt.Sprintf("Updating members for usergroup %s with %d members", usergroupID, len(members)))

	// Join members with comma (Slack API expects comma-separated string)
	membersStr := ""
	if len(members) > 0 {
		for i, member := range members {
			if i > 0 {
				membersStr += ","
			}
			membersStr += member
		}
	}

	_, err := api.UpdateUserGroupMembersContext(ctx, usergroupID, membersStr)
	if err != nil {
		return diag.Errorf("error updating usergroup members: %s", err)
	}

	tflog.Info(ctx, fmt.Sprintf("Usergroup members updated successfully (%d members)", len(members)))
	return nil
}
