mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service.big.aliases == tolist(["service-1", "service-2"])
    error_message = "wrong aliases in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.api_document_path == "api/doc/path.yaml"
    error_message = "wrong api_document_path in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.description == "Scorecard Description"
    error_message = "wrong description in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.framework == "Scorecard Framework"
    error_message = "wrong framework in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.language == "Scorecard Language"
    error_message = "wrong language in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.lifecycle_alias == "alpha"
    error_message = "wrong lifecycle_alias in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.name == "Big Scorecard"
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
    condition     = opslevel_service.big.tags == tolist(["key1:value1"])
    error_message = "wrong tags in opslevel_service.big"
  }

  assert {
    condition     = opslevel_service.big.tier_alias == "Scorecard Tier"
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
    condition     = opslevel_service.small.name == "Small Scorecard"
    error_message = "wrong name in opslevel_service.small"
  }
}
