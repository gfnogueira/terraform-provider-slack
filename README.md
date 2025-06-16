# 🚀 Terraform Slack Provider

A custom Terraform provider to manage **Slack channels**, **users**, and **user groups** with fine control over channel metadata and membership.

---

## ✅ Features

- **Slack Channels**
  - Create public/private channels
  - Set purpose and topic
  - Manage members via Terraform (with drift detection)
  - Rename or reuse existing channels

- **Slack Users**
  - Data source to fetch user by email
  - Data source to list all users

- **Slack Channels (Data Sources)**
  - Get a single channel or list all (with filters)

- **Slack User Groups**
  - Lookup by handle (data source)
  - *(Planned)* Manage user groups (resource)

- **Channel Membership Diff Tool**
  - Output unmanaged channels (not declared in code)

---

## 📦 Installation (Local Dev)

```hcl
provider "slack" {
  token = var.slack_token
}
```

~/.terraformrc:
```hcl
provider_installation {
  dev_overrides {
    "slack/slack" = "/ABSOLUTE/PATH/TO/terraform-provider-slack"
  }
  direct {}
}
```

---

## 🔐 Slack Token Scopes

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

## 🧪 Example

```hcl
resource "slack_channel" "example" {
  name       = "alerts"
  is_private = false
  topic      = "System alerts"
  purpose    = "Terraform-managed"
}
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## 📄 License

[GNU GPLv3](./LICENSE)