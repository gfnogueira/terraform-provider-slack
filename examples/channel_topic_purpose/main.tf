variable "slack_token" {
  description = "Slack API token (xoxb-...)"
  type        = string
  sensitive   = true
}

provider "slack" {
  token = var.slack_token
}

# Public Slack channel with topic and purpose
resource "slack_channel" "public_channel" {
  name       = "tf-public"
  is_private = false
  topic      = "Public channel managed by Terraform"
  purpose    = "Example of a public channel created with Terraform"
}

# Private Slack channel with topic and purpose
resource "slack_channel" "private_channel" {
  name       = "tf-private"
  is_private = true
  topic      = "Private channel via Terraform"
  purpose    = "Example of a private channel created with Terraform"
}