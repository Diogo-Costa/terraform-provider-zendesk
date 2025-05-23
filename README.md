# Terraform Provider for Zendesk OAuth

This Terraform provider allows you to manage OAuth clients and tokens in Zendesk. It provides resources for creating and managing OAuth clients and their associated access tokens.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using `go build -o terraform-provider-zendesk`

## Using the provider

To use the provider, include it in your Terraform configuration:

```hcl
terraform {
  required_providers {
    zendesk = {
      source = "diogocosta/terraform-provider-zendesk"
    }
  }
}

provider "zendesk" {
  subdomain = "your-subdomain"  # The subdomain of your Zendesk account (e.g., company in company.zendesk.com)
  email     = "admin@example.com"  # Your Zendesk admin email
  api_token = "your-api-token"  # Your Zendesk API token
}

# Create an OAuth client
resource "zendesk_oauth_client" "example" {
  name        = "Example Client"
  identifier  = "example_client"
  kind        = "public"
  description = "OAuth client for my custom application"  # Optional description
}

# Create an OAuth token with expiration
resource "zendesk_oauth_token" "example" {
  client_id  = zendesk_oauth_client.example.id
  scopes     = ["read", "write"]
  expires_at = "2024-12-31T23:59:59Z"  # Optional expiration date
}
```

### Environment Variables

You can also use environment variables to configure the provider:

- `ZENDESK_SUBDOMAIN` - The subdomain of your Zendesk account
- `ZENDESK_EMAIL` - Your Zendesk admin email
- `ZENDESK_API_TOKEN` - Your Zendesk API token

## Resources

### `zendesk_oauth_client`

Manages a Zendesk OAuth client.

#### Argument Reference

* `name` - (Required) The name of the OAuth client.
* `identifier` - (Required) The unique identifier of the OAuth client.
* `kind` - (Required) The kind of OAuth client (e.g., 'public').
* `description` - (Optional) A description of the OAuth client.

#### Attribute Reference

* `id` - The ID of the OAuth client.

### `zendesk_oauth_token`

Manages a Zendesk OAuth token.

#### Argument Reference

* `client_id` - (Required) The ID of the OAuth client.
* `scopes` - (Required) The list of scopes granted to the OAuth token.
* `expires_at` - (Optional) The expiration date of the token in ISO 8601 format (e.g., '2024-12-31T23:59:59Z'). If not set, the token will not expire.

#### Attribute Reference

* `id` - The ID of the OAuth token.
* `full_token` - The full OAuth token value (only available after creation).

## Examples

### Basic OAuth Client and Token

```hcl
# Create a basic OAuth client
resource "zendesk_oauth_client" "basic" {
  name       = "Basic Client"
  identifier = "basic_client"
  kind       = "public"
}

# Create a non-expiring token
resource "zendesk_oauth_token" "basic" {
  client_id = zendesk_oauth_client.basic.id
  scopes    = ["read"]
}
```

### OAuth Client with Description and Expiring Token

```hcl
# Create an OAuth client with description
resource "zendesk_oauth_client" "app" {
  name        = "My Custom App"
  identifier  = "my_custom_app"
  kind        = "public"
  description = "OAuth client for my custom application that integrates with Zendesk"
}

# Create a token that expires in one year
resource "zendesk_oauth_token" "app_token" {
  client_id  = zendesk_oauth_client.app.id
  expires_at = timeadd(timestamp(), "8760h")  # 1 year from now
  scopes     = [
    "read",
    "write",
    "tickets:read",
    "tickets:write",
    "users:read",
    "users:write"
  ]
}

# Output the token value (sensitive)
output "oauth_token" {
  value     = zendesk_oauth_token.app_token.full_token
  sensitive = true
}
```

## Contributing

We welcome contributions to this provider! Here's how you can help:

### Development Requirements

- [Go](http://www.golang.org) version 1.21+
- [Terraform](https://www.terraform.io/downloads.html) version 1.0+
- A Zendesk account for testing

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork
3. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Make your changes
5. Run tests:
   ```bash
   go test -v ./...
   ```

### Running Acceptance Tests

To run the acceptance tests, you'll need to set the following environment variables:
```bash
export ZENDESK_SUBDOMAIN="your-subdomain"
export ZENDESK_EMAIL="your-email"
export ZENDESK_API_TOKEN="your-token"
```

Then run the acceptance tests:
```bash
make testacc
```

**Note:** Acceptance tests create real resources in your Zendesk account.

### Submitting Changes

1. Update documentation as needed
2. Add tests for new features
3. Run the test suite
4. Commit your changes (make sure your commit messages are clear)
5. Push to your fork
6. Create a Pull Request

### Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the examples/ directory with examples demonstrating new features
3. The PR will be merged once you have the sign-off of the maintainers

### Release Process

Releases are automatically created when a new tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The release workflow will:
1. Build the provider for all supported platforms
2. Create a draft GitHub release
3. Sign the release with GPG
4. Upload all artifacts

### Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is *required*).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. 