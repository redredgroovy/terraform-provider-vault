package vault

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "URL to the Vault server",
				Required:    true,
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "User ID for app-id authentication",
				Optional:    true,
				Default:     "",
			},
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "App ID for app-id authentication",
				Optional:    true,
				Default:     "",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Token for direct authentication",
				Optional:    true,
				Default:     "",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vault_secret": resourceVaultSecret(),
		},

		ConfigureFunc: providerConfigure,
	}
}

// Perform the initial configuration and authentication for the Vault client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Configuring Vault provider")

	// Copy existing ResourceData set into the vaultConfig struct
	var config vaultConfig
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}

	return config.Client()
}
