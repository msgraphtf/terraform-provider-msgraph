package main

import (
	"context"
	"terraform-provider-msgraph/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name msgraph

func main() {
	providerserver.Serve(context.Background(), msgraph.New, providerserver.ServeOpts{
		Address: "regsistry.terraform.io/msgraphtf/msgraph",
	})

}
