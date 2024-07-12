mock_resource "opslevel_scorecard" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases        = ["scorecard-one", "scorecard-two"]
    categories     = ["Z2lkOi8vlaksjdfsd"]
    passing_checks = 7
    service_count  = 3
    total_checks   = 10
  }
}

