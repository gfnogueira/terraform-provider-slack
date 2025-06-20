---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Slack Provider"
subcategory: "Communication & Messaging"
description: |-
  Terraform provider to manage Slack resources such as channels, users, and user groups.
---

# Slack Provider

The `slack` provider allows you to manage your Slack workspace declaratively using Terraform.  
With this provider, you can automate Slack channel creation, update metadata, sync members, and fetch information about users and user groups.

> ⚠️ **Note**: You need a valid Slack Bot Token (`xoxb-...`) with the necessary scopes to use this provider.

---

## Features

With this provider you can:

- Create, update, and delete public/private channels
- Sync channel membership (declaratively!)
- Manage topic and purpose of channels
- Retrieve Slack users and user groups via data sources
- Import existing Slack channels to manage them via Terraform

---

## Example Usage

```hcl
provider "slack" {
  token = var.slack_token
}

resource "slack_channel" "alerts" {
  name       = "alerts"
  is_private = true
  topic      = "Monitoring alerts"
  purpose    = "Infrastructure team"
}

data "slack_users" "all" {}

output "engineers" {
  value = [for user in data.slack_users.all.users : user.real_name if user.is_admin]
}
```

---

## 🔐 Required Slack OAuth Scopes

Your Slack Bot Token must include the following **OAuth Scopes**, depending on what you plan to manage:

```txt
users:read
users:read.email
channels:read
channels:write
groups:read
groups:write
usergroups:read
```

To obtain this token:

1. Go to https://api.slack.com/apps and **create an app**
2. Under **OAuth & Permissions**, add the scopes listed above
3. Install the app into your workspace
4. Copy the **Bot User OAuth Token** (starts with `xoxb-...`)

---

## Configuration

You can configure the provider by passing the token directly or using an environment variable:

```hcl
provider "slack" {
  token = var.slack_token
}
```

We recommend using an env var for security:

```bash
export SLACK_TOKEN="xoxb-..."
```

---

## Importing Existing Channels

You can import existing Slack channels using their ID:

```bash
terraform import slack_channel.my_channel C0123456789
```

To discover channel IDs, use:

```hcl
data "slack_channels" "all" {}
```

---

## Available Resources

- `slack_channel`: Create and manage Slack channels
- `slack_channel_members`: Manage channel membership (in sync)
- `slack_channel_topic_purpose`: Manage topic/purpose of a channel

## Available Data Sources

- `slack_channel`
- `slack_channels`
- `slack_user`
- `slack_users`
- `slack_usergroup`

---

## Run


```bash
terraform init --upgrade
terraform apply
```

---

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `token` (String, Sensitive) Slack Bot Token (`xoxb-...`)
