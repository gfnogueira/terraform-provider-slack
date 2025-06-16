package test

import (
	"testing"
	"os"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSlackChannel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SLACK_TOKEN"); v == "" {
				t.Skip("SLACK_TOKEN must be set")
			}
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "slack" {
  token = "` + os.Getenv("SLACK_TOKEN") + `"
}

resource "slack_channel" "test" {
  name       = "tf-test-channel"
  is_private = false
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("slack_channel.test", "name", "tf-test-channel"),
					resource.TestCheckResourceAttr("slack_channel.test", "is_private", "false"),
				),
			},
		},
	})
}