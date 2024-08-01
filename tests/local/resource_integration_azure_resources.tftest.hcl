mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_infra_azure_resources_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.last_synced_at == "2024-07-08T13:50:07Z"
    error_message = "wrong last_synced_at for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = can(opslevel_integration_azure_resources.example.id)
    error_message = "expected opslevel_integration_azure_resources to have an ID"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.name == "dev"
    error_message = "wrong name for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.tenant_id == "98765432-9876-9876-9876-987654321098"
    error_message = "wrong tenant_id for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.subscription_id == "01234567-0123-0123-0123-012345678901"
    error_message = "wrong subscription_id for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.client_id == "XXX_CLIENT_ID_XXX"
    error_message = "wrong client_id for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.client_secret == "XXX_CLIENT_SECRET_XXX"
    error_message = "wrong client_secret for opslevel_integration_azure_resources.example"
  }

  assert {
    condition     = opslevel_integration_azure_resources.example.aliases == tolist(["alias1", "alias2"])
    error_message = "wrong aliases for opslevel_integration_azure_resources.example"
  }

}
