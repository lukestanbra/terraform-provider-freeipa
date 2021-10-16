package freeipa

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ipa "github.com/lukestanbra/go-freeipa/freeipa"
)

func TestAccUserBasic(t *testing.T) {
	username := "testuser"
	first_name := "Test"
	last_name := "User"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateUserBasic(username, first_name, last_name),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("freeipa_user.this"),
				),
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*ipa.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "freeipa_user" {
			continue
		}

		username := rs.Primary.ID

		_, err := c.UserDel(
			&ipa.UserDelArgs{},
			&ipa.UserDelOptionalArgs{
				UID: &[]string{username},
			},
		)

		if err == nil {
			return fmt.Errorf("Resource still exists")
		}

		notFoundErr := "NotFound"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}

	}

	return nil
}

func testAccCreateUserBasic(username string, first_name string, last_name string) string {
	return fmt.Sprintf(`
	resource "freeipa_user" "this" {
		username   = "%s"
		first_name = "%s"
		last_name  = "%s"
	}
	`, username, first_name, last_name)
}

func testAccUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No username set")
		}

		return nil
	}
}
