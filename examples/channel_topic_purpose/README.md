# Channel with Topic and Purpose

This example creates a Slack channel and sets a topic and purpose for it.

- `tf-public`: a public Slack channel
- `tf-private`: a private Slack channel

Both channels are configured with `topic` and `purpose`.

---

## How to use

1. Export your Slack token (bot or user token with proper permissions):

```bash
export TF_VAR_slack_token="xoxb-xxxxxxxxxxxx"
```
2. Terraform apply

```sh
terraform init
terraform apply
```