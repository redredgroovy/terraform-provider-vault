package vault

import (
	"fmt"
	"log"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

type vaultConfig struct {
	Address string `mapstructure:"address"`
	UserId  string `mapstructure:"user_id"`
	AppId   string `mapstructure:"app_id"`
	Token   string `mapstructure:"token"`
}

// Create a Vault client and authenticate for a token, if necessary
func (c *vaultConfig) Client() (*vaultapi.Client, error) {
	config := vaultapi.DefaultConfig()

	if c.Address != "" {
		config.Address = c.Address
	}
	log.Printf("[INFO] Vault client configured with address: %s", config.Address)

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Set the authentication token, if provided
	if c.Token != "" {

		log.Printf("[INFO] Vault client using token authentication")
		client.SetToken(c.Token)

	} else if c.AppId != "" && c.UserId != "" { // AppId authentication

		log.Printf("[INFO] Vault client using app-id authentication")

		// Build the request JSON body
		body := map[string]interface{}{
			"app_id":  strings.TrimSpace(c.AppId),
			"user_id": strings.TrimSpace(c.UserId),
		}
		secret, err := client.Logical().Write("auth/app-id/login", body)
		if err != nil {
			return nil, err
		}

		// Set the token if authentication was successful
		client.SetToken(secret.Auth.ClientToken)

	} else { // No authentication provided

		return nil, fmt.Errorf("Vault provider requires either 'token' or 'app_id + user_id'")

	}

	return client, nil
}
