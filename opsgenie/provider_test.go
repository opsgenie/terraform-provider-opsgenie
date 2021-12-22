package opsgenie

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var providerFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	providerFactories = map[string]func() (*schema.Provider, error){
		"opsgenie": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}
