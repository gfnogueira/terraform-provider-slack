provider "slack" {
  token = var.slack_token
}

data "slack_users_group" "example" {
  handle = "developers"
}
