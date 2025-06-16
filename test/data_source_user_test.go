package test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSlackUser(t *testing.T) {
	username := os.Getenv("SLACK_TEST_USER")
	if username == "" {
		t.Skip("SLACK_TEST_USER must be set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("SLACK_TOKEN") == "" {
				t.Fatal("SLACK_TOKEN must be set")
			}
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "slack" {
  token = "` + os.Getenv("SLACK_TOKEN") + `"
}

data "slack_user" "example" {
  name = "` + username + `"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.slack_user.example", "id"),
					resource.TestCheckResourceAttr("data.slack_user.example", "name", username),
				),
			},
		},
	})
}