package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"pihole": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PIHOLE_URL"); v == "" {
		t.Fatal("PIHOLE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("PIHOLE_PASSWORD"); v == "" {
		t.Fatal("PIHOLE_PASSWORD must be set for acceptance tests")
	}
}
