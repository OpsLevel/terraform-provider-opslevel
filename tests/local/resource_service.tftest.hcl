mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = alltrue([
      contains(opslevel_service.big.aliases, "service-1"),
      contains(opslevel_service.big.aliases, "service-2"),
    ])
    error_message = "wrong aliases in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.api_document_path == "api/doc/path.yaml"
    error_message = "wrong api_document_path in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.description == "Service Description"
    error_message = "wrong description in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.framework == "Service Framework"
    error_message = "wrong framework in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.language == "Service Language"
    error_message = "wrong language in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.lifecycle_alias == "alpha"
    error_message = "wrong lifecycle_alias in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.name == "Big Service"
    error_message = "wrong name in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.owner == "team-alias"
    error_message = "wrong owner in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.preferred_api_document_source == "PULL"
    error_message = "wrong preferred_api_document_source in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.product == "Mock Product"
    error_message = "wrong product in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.tags == tolist(["key1:value1", "key2:value2"])
    error_message = "wrong tags in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.tier_alias == "Service Tier"
    error_message = "wrong tier_alias in opslevel_service.big"
  }

}

run "resource_service_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_service.small.id)
    error_message = "id attribute missing from filter in opslevel_rubric_level.small"
  }

  assert {
    condition     = opslevel_service.small.name == "Small Service"
    error_message = "wrong name in opslevel_service.small"
  }
}
