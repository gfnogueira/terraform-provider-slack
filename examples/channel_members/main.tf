provider "slack" {
  token = var.slack_token
}

resource "slack_channel" "members_example" {
  name       = "members-channel"
  is_private = false
  members    = ["U12345678", "U87654321"]
}