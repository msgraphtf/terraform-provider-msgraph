package main

import (
	"context"
	"terraform-provider-msgraph/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), msgraph.New, providerserver.ServeOpts{
		Address: "regsistry.terraform.io/hsheppard/msgraph",
	})

}
