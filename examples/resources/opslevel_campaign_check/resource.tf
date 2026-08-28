data "opslevel_team" "platform" {
  alias = "platform"
}

resource "opslevel_check_tool_usage" "arize" {
  name     = "Arize tool usage"
  category = data.opslevel_rubric_category.reliability.id
  level    = data.opslevel_rubric_level.bronze.id

  tool_category = "observability"
}

resource "opslevel_campaign" "upgrade_rails" {
  name     = "Upgrade to Rails 7"
  owner_id = data.opslevel_team.platform.id
}

# Copies the rubric check onto the campaign. The copy is a new, independent
# check - later edits to opslevel_check_tool_usage.arize do not change it.
resource "opslevel_campaign_check" "arize" {
  campaign_id          = opslevel_campaign.upgrade_rails.id
  copied_from_check_id = opslevel_check_tool_usage.arize.id

  owner   = data.opslevel_team.platform.id
  notes   = "Contact #platform if this check is failing."
  enabled = true
}

# Copy several checks onto one campaign. A campaign can hold at most 25.
resource "opslevel_campaign_check" "rails_checks" {
  for_each = {
    gemfile = opslevel_check_repository_file.gemfile.id
    ci      = opslevel_check_repository_file.ci_config.id
  }

  campaign_id          = opslevel_campaign.upgrade_rails.id
  copied_from_check_id = each.value
}
