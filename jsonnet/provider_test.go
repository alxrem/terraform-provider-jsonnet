package jsonnet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
