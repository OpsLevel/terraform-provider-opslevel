mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_scorecard_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.affects_overall_service_levels == true
    error_message = "wrong id in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.aliases == tolist(["sc-alias-one", "sc-alias-two"])
    error_message = "wrong aliases in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.description == "mock-scorecard-description"
    error_message = "wrong description in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.filter_id == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS84Njcw"
    error_message = "wrong filter_id in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.id == "Z2lkOi8vb3BzbGV2ZWwvU2VybqijZS84Npic"
    error_message = "wrong id in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.name == "mock-scorecard-name"
    error_message = "wrong name in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.owner_id == "Z2lkOi8vb3BzbGV2ZWwvU2VybqijZS84Noqp"
    error_message = "wrong owner_id in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.passing_checks == 20
    error_message = "wrong passing_checks in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.service_count == 10
    error_message = "wrong service_count in opslevel_scorecard mock"
  }

  assert {
    condition     = data.opslevel_scorecard.mock_scorecard.total_checks == 50
    error_message = "wrong total_checks in opslevel_scorecard mock"
  }

}
