variables {
  campaign_name     = "TF Test Campaign Check"
  owner_id          = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8x" # replace with a valid team ID in your org
  source_check_name = "TF Test Source Check"
}

run "resource_campaign_check_create" {
  variables {
    campaign_name     = var.campaign_name
    owner_id          = var.owner_id
    source_check_name = var.source_check_name
  }

  module {
    source = "./campaign_check"
  }

  assert {
    condition = alltrue([
      can(opslevel_campaign_check.test.id),
      can(opslevel_campaign_check.test.type),
      can(opslevel_campaign_check.test.name),
    ])
    error_message = "expected campaign check to have id, type, and name"
  }

  # The whole point of the resource: the copy is a new, independent check, so its
  # id must differ from the rubric check it was copied from.
  assert {
    condition     = opslevel_campaign_check.test.id != opslevel_check_manual.source.id
    error_message = "the copy must have its own id, distinct from the check it was copied from"
  }

  # sourceCheck is what makes the copy traceable back to its origin.
  assert {
    condition     = opslevel_campaign_check.test.source_check_id == opslevel_check_manual.source.id
    error_message = "source_check_id should report the rubric check the copy was made from"
  }

  assert {
    condition     = opslevel_campaign_check.test.campaign_id == opslevel_campaign.test.id
    error_message = "campaign check should belong to the campaign it was copied onto"
  }

  # The copy inherits its name from the source rather than being set in config.
  assert {
    condition     = opslevel_campaign_check.test.name == var.source_check_name
    error_message = "the copy should inherit the source check's name"
  }

  # Copies are always created disabled by the backend.
  assert {
    condition     = opslevel_campaign_check.test.enabled == false
    error_message = "a copy should default to disabled"
  }

  # Requirement: a campaign's checks must be visible from Terraform.
  assert {
    condition     = length(data.opslevel_campaign.test.checks) == 1
    error_message = "the campaign data source should report the copied check"
  }

  assert {
    condition     = data.opslevel_campaign.test.checks[0].id == opslevel_campaign_check.test.id
    error_message = "the campaign data source should report the copy's own id"
  }

  assert {
    condition     = data.opslevel_campaign.test.checks[0].source_check_id == opslevel_check_manual.source.id
    error_message = "the campaign data source should report where the copy came from"
  }
}

# The copy is updated in place rather than being deleted and re-copied, which
# would discard its results and history.
run "resource_campaign_check_update_in_place" {
  variables {
    campaign_name     = var.campaign_name
    owner_id          = var.owner_id
    source_check_name = var.source_check_name
    notes             = "Updated by the acceptance test"
    enabled           = true
  }

  module {
    source = "./campaign_check"
  }

  assert {
    condition     = opslevel_campaign_check.test.notes == "Updated by the acceptance test"
    error_message = "notes should be updatable on a campaign check"
  }

  assert {
    condition     = opslevel_campaign_check.test.enabled == true
    error_message = "enabled should be updatable on a campaign check"
  }
}
