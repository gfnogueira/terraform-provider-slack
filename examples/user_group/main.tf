provider "slack" {
  token = var.slack_token
}

data "slack_users_group" "engineering" {
  handle = "engineering"
}