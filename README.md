# 🚀 terraform-provider-slack

A custom Terraform provider for managing Slack resources such as channels, user groups, and users.

## ⚙️ Features

This provider allows you to interact with the Slack API to manage:

- [x] **Slack Channels**  
  - Create public or private channels
  - Rename existing channels
  - Manage channel members
  - Set channel **purpose** and **topic**
  - Automatically reuses existing channels (including archived ones)
- [x] **Slack Users** *(Data Source)*  
  - Fetch user information by email or ID
- [x] **Slack Channels** *(Data Sources)*  
  - Fetch channel by name or ID
  - Fetch all channels (list)
- [x] **Slack User Groups** *(Data Source)*  
  - Fetch user group by handle or name
- [x] **Slack User Groups** *(Resource)*  
  > ⚠️ **Note**: Slack user group management is only available on Slack paid plans.
  - (Planned) Create and update user groups, including name, handle, description, and users

---

## 🧪 Example Usage

### Provider Setup

```
provider "slack" {
  token = var.slack_token  # or use the SLACK_TOKEN environment variable
}
```

### Create a Slack Channel

```
resource "slack_channel" "devops" {
  name       = "devops-alerts"
  is_private = false

  members = [
    "U1234567890", # user IDs
    "U0987654321",
  ]

  purpose = "Channel for DevOps alerts and automation"
  topic   = "Alerts, Deployments and Monitoring discussions"
}
```

---

## 📦 Install the Provider

Build from source:

```
go build -o terraform-provider-slack
```

Then configure the plugin in your Terraform CLI config:

```
# ~/.terraformrc or %APPDATA%\terraform.rc

provider_installation {
  dev_overrides {
    "slack/slack" = "/path/to/compiled/provider"
  }
  direct {}
}
```

---

## 🔐 Authentication

This provider requires a **Slack token** with the following scopes:

```
users:read
users:read.email
channels:read
channels:write
groups:read
groups:write
usergroups:read
```

You can provide it via:

- The `token` argument in the provider block
- Or the `SLACK_TOKEN` environment variable

---

## 🔍 Planned Features

- [ ] `slack_usergroup` resource (create/update user groups)
- [ ] Channel archiving support
- [ ] Slack bot lifecycle testing
- [ ] Slack workflows or scheduled messages (stretch goal)

---

## 🧠 Development Notes

- Built with [terraform-plugin-sdk v2](https://github.com/hashicorp/terraform-plugin-sdk)
- Uses [slack-go/slack](https://github.com/slack-go/slack)
- Tested on personal and team Slack workspaces

---

## 🤝 Contributing

Pull requests are welcome! If you have ideas or need support for specific resources (like bots or apps), feel free to open an issue.

---

## 📄 License

MIT — see [LICENSE](./LICENSE) for details.
