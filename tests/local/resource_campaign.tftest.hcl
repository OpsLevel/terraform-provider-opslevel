mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_campaign_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_campaign.big.name == "Big Campaign"
    error_message = "wrong name in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.owner_id == var.test_id
    error_message = "wrong owner_id in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.filter_id == var.test_id
    error_message = "wrong filter_id in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.project_brief == "This is a big campaign"
    error_message = "wrong project_brief in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.start_date == "2026-07-01"
    error_message = "wrong start_date in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.target_date == "2026-09-30"
    error_message = "wrong target_date in opslevel_campaign.big"
  }

  assert {
    condition     = can(opslevel_campaign.big.id)
    error_message = "id attribute missing from opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.status == "draft"
    error_message = "wrong status in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.html_url == "https://app.opslevel.com/campaigns/test"
    error_message = "wrong html_url in opslevel_campaign.big"
  }

  assert {
    condition     = length(opslevel_campaign.big.check_ids) == 1
    error_message = "wrong number of check_ids in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.check_ids[0] == var.test_id
    error_message = "wrong check_ids[0] in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.reminder.frequency_unit == "week"
    error_message = "wrong reminder.frequency_unit in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.reminder.frequency == 1
    error_message = "wrong reminder.frequency in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.reminder.time_of_day == "09:30"
    error_message = "wrong reminder.time_of_day in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.reminder.timezone == "America/Chicago"
    error_message = "wrong reminder.timezone in opslevel_campaign.big"
  }

  assert {
    condition     = length(opslevel_campaign.big.reminder.channels) == 2
    error_message = "wrong number of reminder.channels in opslevel_campaign.big"
  }

  assert {
    condition     = length(opslevel_campaign.big.reminder.days_of_week) == 2
    error_message = "wrong number of reminder.days_of_week in opslevel_campaign.big"
  }

  assert {
    condition     = opslevel_campaign.big.reminder.default_slack_channel == "#platform-eng"
    error_message = "wrong reminder.default_slack_channel in opslevel_campaign.big"
  }
}

run "resource_campaign_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_campaign.small.name == "Small Campaign"
    error_message = "wrong name in opslevel_campaign.small"
  }

  assert {
    condition     = opslevel_campaign.small.owner_id == var.test_id
    error_message = "wrong owner_id in opslevel_campaign.small"
  }

  assert {
    condition     = can(opslevel_campaign.small.id)
    error_message = "id attribute missing from opslevel_campaign.small"
  }
}
