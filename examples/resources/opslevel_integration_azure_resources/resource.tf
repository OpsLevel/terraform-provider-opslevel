resource "opslevel_integration_azure_resources" "dev" {
  client_id       = "XXX_CLIENT_ID_XXX"
  client_secret   = "XXX_CLIENT_SECRET_XXX"
  name            = "Azure Integration"
  subscription_id = "01234567-0123-0123-0123-012345678901"
  tenant_id       = "98765432-9876-9876-9876-987654321098"

  ownership_tag_keys      = ["owner", "team", "group"]
  ownership_tag_overrides = true
}
