terraform {
  required_providers {
    zendesk = {
      source = "diogocosta/terraform-provider-zendesk"
    }
  }
}

provider "zendesk" {
  # Configure using environment variables:
  # ZENDESK_SUBDOMAIN
  # ZENDESK_EMAIL
  # ZENDESK_API_TOKEN
}

# Example 1: Basic OAuth client and token
resource "zendesk_oauth_client" "basic" {
  name       = "Basic Client"
  identifier = "basic_client"
  kind       = "public"
}

resource "zendesk_oauth_token" "basic" {
  client_id = zendesk_oauth_client.basic.id
  scopes    = ["read"]
}

# Example 2: OAuth client with description and expiring token
resource "zendesk_oauth_client" "app" {
  name        = "My Custom App"
  identifier  = "my_custom_app"
  kind        = "public"
  description = "OAuth client for my custom application that integrates with Zendesk"
}

resource "zendesk_oauth_token" "app_token" {
  client_id  = zendesk_oauth_client.app.id
  expires_at = timeadd(timestamp(), "8760h") # 1 year from now
  scopes     = [
    "read",
    "write",
    "tickets:read",
    "tickets:write",
    "users:read",
    "users:write"
  ]
}

# Example 3: OAuth client with specific expiration date
resource "zendesk_oauth_client" "scheduled" {
  name        = "Scheduled Access App"
  identifier  = "scheduled_app"
  kind        = "public"
  description = "OAuth client with scheduled access period"
}

resource "zendesk_oauth_token" "scheduled_token" {
  client_id  = zendesk_oauth_client.scheduled.id
  expires_at = "2024-12-31T23:59:59Z" # Specific expiration date
  scopes     = [
    "read",
    "tickets:read",
    "users:read"
  ]
}

# Outputs
output "basic_token" {
  value     = zendesk_oauth_token.basic.full_token
  sensitive = true
}

output "app_token" {
  value     = zendesk_oauth_token.app_token.full_token
  sensitive = true
}

output "scheduled_token" {
  value     = zendesk_oauth_token.scheduled_token.full_token
  sensitive = true
} 