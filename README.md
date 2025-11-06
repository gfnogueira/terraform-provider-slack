# Terraform Provider for Slack

[![License](https://img.shields.io/github/license/gfnogueira/terraform-provider-slack)](LICENSE)
[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-Slack%20Provider-brightgreen?logo=terraform)](https://registry.terraform.io/providers/gfnogueira/slack/latest)

A Terraform provider to manage Slack channels, users, user groups, and more — directly from your infrastructure code.

---

## ⚙️ Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.3.0
- Go >= 1.22 (only for local development)

---

## 📦 Installation

To use this provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    slack = {
      source  = "gfnogueira/slack"
      version = "0.1.5"
    }
  }
}
```

Then run:

```bash
terraform init
```

---

## 🔐 Authentication

The provider uses a Slack bot token (with appropriate scopes).  
Add the token to your environment:

```bash
export SLACK_TOKEN="xoxb-1234567890-..."
```

Or use it in the provider block:

```hcl
provider "slack" {
  token = var.slack_token
}
```

---

## 🚀 Features

This release supports:

### Resources

- `slack_channel`
  - Create public/private channels
  - Update name, topic, and purpose
  - Manage members
  - Delete channels

### Data Sources

- `slack_channel`
- `slack_channels`
- `slack_user`
- `slack_users`
- `slack_users_group`

---

## 📂 Examples

See the [`examples/`](./examples) directory for complete use cases:

- Create a channel with members
- Set channel topic and purpose
- Query users and user groups

---

## 📄 Documentation

Documentation for all resources and data sources is available under the [`docs/`](./docs) folder and on the Terraform Registry (once published).

---

## 🛠 Development

```bash
make build
make test
```

---

## 🤝 Contributing

Please see [`CONTRIBUTING.md`](./CONTRIBUTING.md).

---

## 🛡 License

This project is licensed under the [GNU GPLv3](./LICENSE).

---

## 📬 Author

Created and maintained by [@gfnogueira](https://github.com/gfnogueira).
