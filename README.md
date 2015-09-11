# terraform-provider-vault
[Terraform](https://www.terraform.io/) provider for reading secrets from [Hashicorp Vault](https://www.vaultproject.io/)

This provider allows you to retrieve secrets from a Vault instance and access them as variables within your Terraform configurations. It can provide convenient access to sensitive material, including:
* AWS credentials
* SSH keys
* SSL certificate keys

## Important Security Notice
This provider currently stores **unencrypted** secrets as plaintext in your `.tfstate` files! Do **not** use this module unless you have taken steps to ensure your `.tfstate` files are encrypted or otherwise secured. This provider should be considered unsafe until a solution has been found for [terraform:#516](https://github.com/hashicorp/terraform/issues/516).

## Usage

### Provider configuration
This module implements a `vault` provider.
```apache
provider "vault" {
    address = "https://vault:8200"
    app_id = "..."
    user_id = "..."
    token = "..."
}
```

##### Arguments

The Vault provider currently supports `token`, `app-id`, `userpass`, and `ldap`  authentication.

* `address` - (required) URL to the Vault API
* `token` - Explicit token for `token` authentication
* `app_id` - Application ID for `app-id` authentication
* `user_id` - User ID for `app-id` authentication
* `user` - Username for `userpass` authentication
* `pass` - Password for `userpass` authentication
* `ldapuser` - Username for `ldap` authentication
* `ldappass` - Password for `ldap` authentication

### Resource configuration
This module implements a `vault_secret` resource.
```apache
resource "vault_secret" "aws" {
    path = "/secret/aws/credentials"
}
```

##### Arguments
* `path` - The path to the secret store for this resource

The contents of `path` will be stored as a map in the `data` variable of the resource.

### Example
```apache
provider "vault" {
    address = "https://vault:8200"
    app_id = "eea928cc-2e83-4db7-8ad2-b90b7bd43542"
    user_id = "${file("~/.user-id")}"
}

resource "vault_secret" "aws" {
    path = "/secret/aws/credentials"
}

# Assuming a Vault entry with the following fields:
#   access_key
#   secret_key
provider "aws" {
    access_key = "${vault_secret.aws.data.access_key}"
    secret_key = "${vault_secret.aws.data.secret_key}"
}

resource "vault_secret" "cert" {
    path = "/secret/certs/www"
}

# Assuming a Vault entry with the following fields:
#   cert
#   key
resource "aws_iam_server_certificate" "www" {
    name = "www"
    certificate_body = "${vault_secret.www.data.cert}"
    private_key = "${vault_secret.www.data.key}"
}
```

## Author

[Derek Moore](https://github.com/redredgroovy)
