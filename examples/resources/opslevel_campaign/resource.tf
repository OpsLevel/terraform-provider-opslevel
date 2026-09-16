data "opslevel_team" "platform" {
  alias = "platform"
}

data "opslevel_filter" "tier_1" {
  filter {
    field = "name"
    value = "Tier 1 Services"
  }
}

resource "opslevel_campaign" "upgrade_rails" {
  name      = "Upgrade to Rails 7"
  owner_id  = data.opslevel_team.platform.id
  filter_id = data.opslevel_filter.tier_1.id

  start_date  = "2026-07-01"
  target_date = "2026-09-30"

  project_brief = <<-EOT
    ## Overview
    All Rails services must upgrade to Rails 7 by end of Q3.

    ## What you need to do
    1. Update your Gemfile to target Rails 7
    2. Run the Rails upgrade checklist
    3. Verify all tests pass
  EOT
}
