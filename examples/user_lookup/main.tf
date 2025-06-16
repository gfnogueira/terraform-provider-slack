provider "slack" {
  token = var.slack_token
}

data "slack_user" "by_email" {
  email = "user@example.com"
}