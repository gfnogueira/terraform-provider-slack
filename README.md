# 🚀 terraform-provider-slack

A custom Terraform provider for managing Slack resources such as channels, users, and user groups with advanced features.

---

## ⚙️ Features

### ✅ Slack Channels (Resource: `slack_channel`)
- Create public or private channels
- Reuse existing channels (including archived ones)
- Manage members (with automatic drift detection)
- Set channel **purpose** and **topic**
- Rename channels via Terraform
- Detect and warn about extra members not declared in code

### ✅ Slack Users
- `data "slack_user"`: Fetch Slack user by email
- `data "slack_users"`: Fetch all users, optionally filter by domain

### ✅ Slack Channels (Data Sources)
- `data "slack_channel"`: Fetch info for a specific channel by name
- `data "slack_channels"`: List all channels with filters (prefix, type, archived)

### ✅ Slack User Groups (Paid Slack)
- `data "slack_usergroup"`: Fetch user group by handle
- *(Planned)* `resource "slack_usergroup"`: Create/update user groups (name, handle, description, users)

### ✅ Channel Membership Diff Tool
- Detect unmanaged Slack channels (not declared in Terraform)
- Show them via output for visibility without failing the plan

---

## 🧪 Example Usage

### Provider Configuration

```hcl
provider "slack" {
  token = var.slack_token
}
```

### Creating a Channel with Members

```hcl
resource "slack_channel" "devops" {
  name       = "devops-alerts"
  is_private = false

  members = [
    data.slack_user.alice.id,
    data.slack_user.bob.id,
  ]

  purpose = "Channel for DevOps alerts and automation"
  topic   = "Deployments and monitoring discussions"
}
```

### Fetching All Users

```hcl
data "slack_users" "all" {}

output "user_emails" {
  value = [for u in data.slack_users.all.users : u.email]
}
```

### Detecting Unmanaged Channels

```hcl
data "slack_channels" "all" {}

locals {
  defined_channels = toset([for c in slack_channel.managed : c.name])
  all_channels     = toset([for c in data.slack_channels.all.channels : c.name])
  unmanaged        = setsubtract(local.all_channels, local.defined_channels)
}

output "unmanaged_slack_channels" {
  value       = local.unmanaged
  description = "Channels not managed by Terraform"
}
```

---

## 📦 Install the Provider (Development)

Build and link locally:

```sh
go build -o terraform-provider-slack

# ~/.terraformrc
provider_installation {
  dev_overrides {
    "slack/slack" = "/ABSOLUTE/PATH/TO/terraform-provider-slack"
  }
  direct {}
}
```

---

## 🔐 Required Slack Token Scopes

Make sure your Slack token has these scopes:

```text
users:read
users:read.email
channels:read
channels:write
groups:read
groups:write
usergroups:read
```

---

## 🧠 Tech Stack

- Built with [terraform-plugin-sdk v2](https://github.com/hashicorp/terraform-plugin-sdk)
- Uses [slack-go/slack](https://github.com/slack-go/slack)
- Compatible with Slack Free and Paid workspaces

---

## 🚧 Roadmap

- [ ] `slack_usergroup` resource
- [ ] Channel archiving lifecycle support
- [ ] Output user names instead of IDs in diffs
- [ ] Slack workflow/task automation (exploration)

---

## 🤝 Contributing

Pull requests are welcome! Open issues for bugs or feature requests.

---

## 📄 License

MIT — see [LICENSE](./LICENSE)