provider "slack" {
  token = var.slack_token
}

resource "slack_channel" "example" {
  name       = "terraform-example"
  is_private = false
}