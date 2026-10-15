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

  # Send a weekly Slack and email reminder to component owners while the
  # campaign is in progress. days_of_week is only valid with a weekly cadence.
  reminder = {
    channels       = ["slack", "email"]
    frequency      = 1
    frequency_unit = "week"
    days_of_week   = ["monday", "thursday"]
    time_of_day    = "09:30"
    timezone       = "America/Chicago"
    message        = "Friendly reminder to complete your Rails 7 upgrade checks."

    # Fallback channel used when a team has no default Slack contact.
    default_slack_channel = "#platform-eng"
  }
}
