package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	vaultapi "github.com/hashicorp/vault/api"
)

func resourceVaultSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "path to the secret store",
				Required:    true,
			},
			"data": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "contents of the secret storage at 'path'",
				Computed:    true,
			},
		},
	}
}

// No-op: Vault provider is currently read-only
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

// Reads a stored data set from the Vault server
func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	vault := m.(*vaultapi.Client)

	path := d.Get("path").(string)

	// Build a simple resource ID
	d.SetId(fmt.Sprintf("path:%s", path)) // TODO: make this more unique

	secret, err := vault.Logical().Read(path)
	if err != nil {
		return err
	}

	// Set 'data' to the entire content of the Data map
	data := make(map[string]string)

	// A 404 returns nil, but it is not reported as an error by the API
	if secret != nil {
		for k, v := range secret.Data {
			data[k] = resourceDecode(v.(string))
		}
	}

	d.Set("data", data)
	return nil
}

// Test for base64 value and decode if found
func resourceDecode(str string) string {
	prefix := "base64:"
	if strings.HasPrefix(str, prefix) {
		data, _ := base64.StdEncoding.DecodeString(str[len(prefix):])
		str = string(data)
	}
	return str
}

// No-op: Vault provider is currently read-only
func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

// No-op: Vault provider is currently read-only
func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
