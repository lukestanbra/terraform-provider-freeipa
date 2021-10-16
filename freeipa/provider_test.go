package freeipa

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider("dev")()
	testAccProviders = map[string]*schema.Provider{
		"freeipa": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider("dev")()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("FREEIPA_USERNAME"); err == "" {
		t.Fatal("FREEIPA_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("FREEIPA_PASSWORD"); err == "" {
		t.Fatal("FREEIPA_PASSWORD must be set for acceptance tests")
	}
	if err := os.Getenv("FREEIPA_HOST"); err == "" {
		t.Fatal("FREEIPA_HOST must be set for acceptance tests")
	}
}
