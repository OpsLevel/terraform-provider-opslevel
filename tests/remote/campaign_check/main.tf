data "opslevel_rubric_categories" "all" {}

data "opslevel_rubric_levels" "all" {}

# The rubric check that gets copied onto the campaign. Levels are ordered by
# index, and index 0 cannot hold checks, so the second level is used.
resource "opslevel_check_manual" "source" {
  name     = var.source_check_name
  category = data.opslevel_rubric_categories.all.rubric_categories[0].id
  level    = data.opslevel_rubric_levels.all.rubric_levels[1].id

  update_frequency = {
    starting_date = "2026-09-01T00:00:00.000Z"
    time_scale    = "week"
    value         = 1
  }
  update_requires_comment = false
}

resource "opslevel_campaign" "test" {
  name     = var.campaign_name
  owner_id = var.owner_id
}

resource "opslevel_campaign_check" "test" {
  campaign_id          = opslevel_campaign.test.id
  copied_from_check_id = opslevel_check_manual.source.id

  notes   = var.notes
  enabled = var.enabled
}

data "opslevel_campaign" "test" {
  identifier = opslevel_campaign.test.id

  depends_on = [opslevel_campaign_check.test]
}
