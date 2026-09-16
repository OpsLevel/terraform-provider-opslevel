variables {
  name          = "TF Test Campaign"
  owner_id      = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8x" # replace with a valid team ID in your org
  project_brief = "Integration test campaign created by Terraform"
}

run "resource_campaign_create_draft" {
  variables {
    name          = var.name
    owner_id      = var.owner_id
    project_brief = var.project_brief
  }

  module {
    source = "./campaign"
  }

  assert {
    condition = alltrue([
      can(opslevel_campaign.test.id),
      can(opslevel_campaign.test.html_url),
      can(opslevel_campaign.test.status),
    ])
    error_message = "expected campaign to have id, html_url, and status"
  }

  assert {
    condition     = opslevel_campaign.test.name == var.name
    error_message = "campaign name does not match"
  }

  assert {
    condition     = opslevel_campaign.test.status == "draft"
    error_message = "new campaign should be in draft status"
  }

  assert {
    condition     = opslevel_campaign.test.project_brief == var.project_brief
    error_message = "campaign project_brief does not match"
  }
}

run "resource_campaign_schedule" {
  variables {
    name          = var.name
    owner_id      = var.owner_id
    project_brief = var.project_brief
    start_date    = "2026-08-01"
    target_date   = "2026-12-31"
  }

  module {
    source = "./campaign"
  }

  assert {
    condition     = opslevel_campaign.test.start_date == "2026-08-01"
    error_message = "campaign start_date does not match"
  }

  assert {
    condition     = opslevel_campaign.test.target_date == "2026-12-31"
    error_message = "campaign target_date does not match"
  }

  assert {
    condition     = opslevel_campaign.test.status == "scheduled"
    error_message = "campaign should be in scheduled status after setting dates"
  }
}

run "resource_campaign_unschedule" {
  variables {
    name          = "TF Test Campaign Updated"
    owner_id      = var.owner_id
    project_brief = "Updated project brief"
  }

  module {
    source = "./campaign"
  }

  assert {
    condition     = opslevel_campaign.test.name == "TF Test Campaign Updated"
    error_message = "campaign name was not updated"
  }

  assert {
    condition     = opslevel_campaign.test.project_brief == "Updated project brief"
    error_message = "campaign project_brief was not updated"
  }

  assert {
    condition     = opslevel_campaign.test.status == "draft"
    error_message = "campaign should be back in draft status after removing dates"
  }
}

run "resource_campaign_create_with_reminder" {
  variables {
    name          = var.name
    owner_id      = var.owner_id
    project_brief = var.project_brief
    reminder = {
      channels       = ["slack", "email"]
      frequency      = 1
      frequency_unit = "week"
      days_of_week   = ["monday", "thursday"]
      time_of_day    = "09:30"
      timezone       = "America/Chicago"
      message        = "Please complete your campaign checks."
    }
  }

  module {
    source = "./campaign"
  }

  assert {
    condition     = opslevel_campaign.test.reminder.frequency_unit == "week"
    error_message = "reminder frequency_unit does not match"
  }

  assert {
    condition     = length(opslevel_campaign.test.reminder.channels) == 2
    error_message = "reminder channels do not match"
  }

  assert {
    condition     = length(opslevel_campaign.test.reminder.days_of_week) == 2
    error_message = "reminder days_of_week do not match"
  }

  assert {
    condition     = opslevel_campaign.test.reminder.time_of_day == "09:30"
    error_message = "reminder time_of_day does not match"
  }
}

run "resource_campaign_update_reminder_to_daily" {
  variables {
    name          = var.name
    owner_id      = var.owner_id
    project_brief = var.project_brief
    reminder = {
      channels       = ["slack"]
      frequency      = 2
      frequency_unit = "day"
      time_of_day    = "14:00"
      timezone       = "America/Chicago"
    }
  }

  module {
    source = "./campaign"
  }

  assert {
    condition     = opslevel_campaign.test.reminder.frequency_unit == "day"
    error_message = "reminder frequency_unit was not updated to day"
  }

  assert {
    condition     = opslevel_campaign.test.reminder.frequency == 2
    error_message = "reminder frequency was not updated"
  }

  assert {
    condition     = opslevel_campaign.test.reminder.days_of_week == null
    error_message = "reminder days_of_week should be null for a daily cadence"
  }
}

run "resource_campaign_clear_reminder" {
  variables {
    name          = var.name
    owner_id      = var.owner_id
    project_brief = var.project_brief
  }

  module {
    source = "./campaign"
  }

  assert {
    condition     = opslevel_campaign.test.reminder == null
    error_message = "reminder should be cleared when the block is removed"
  }
}
