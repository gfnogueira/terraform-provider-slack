package test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProviderInit(t *testing.T) {
	if os.Getenv("SLACK_TOKEN") == "" {
		t.Skip("SLACK_TOKEN must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SLACK_TOKEN"); v == "" {
				t.Fatal("SLACK_TOKEN must be set for acceptance tests")
			}
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "slack" {
  token = "` + os.Getenv("SLACK_TOKEN") + `"
}
`,
			},
		},
	})
}