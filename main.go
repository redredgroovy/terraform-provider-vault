package main

import (
	"github.com/hashicorp/terraform/terraform"
    "github.com/redredgroovy/terraform-provider-vault"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vault.Provider()
	})
}
