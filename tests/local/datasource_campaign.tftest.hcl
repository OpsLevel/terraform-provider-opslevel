mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_campaign_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.filter_id == "Z2lkOi8vb3BzbGV2ZWwvRmlsdGVyLzEyMw"
    error_message = "wrong filter_id in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.html_url == "https://app.opslevel.com/campaigns/mock"
    error_message = "wrong html_url in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.id == "Z2lkOi8vb3BzbGV2ZWwvQ2FtcGFpZ24vMTIz"
    error_message = "wrong id in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.name == "mock-campaign-name"
    error_message = "wrong name in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.owner_id == "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xMjM"
    error_message = "wrong owner_id in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.project_brief == "mock-project-brief"
    error_message = "wrong project_brief in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.start_date == "2026-07-01"
    error_message = "wrong start_date in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.status == "scheduled"
    error_message = "wrong status in opslevel_campaign mock"
  }

  assert {
    condition     = data.opslevel_campaign.mock_campaign.target_date == "2026-09-30"
    error_message = "wrong target_date in opslevel_campaign mock"
  }
}
