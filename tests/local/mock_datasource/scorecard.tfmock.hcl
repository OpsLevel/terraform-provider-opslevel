mock_data "opslevel_scorecard" {
  defaults = {
    affects_overall_service_levels = true
    aliases                        = ["sc-alias-one", "sc-alias-two"]
    description                    = "mock-scorecard-description"
    filter_id                      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS84Njcw"
    id                             = "Z2lkOi8vb3BzbGV2ZWwvU2VybqijZS84Npic"
    name                           = "mock-scorecard-name"
    owner_id                       = "Z2lkOi8vb3BzbGV2ZWwvU2VybqijZS84Noqp"
    passing_checks                 = 20
    service_count                  = 10
    total_checks                   = 50
  }
}
