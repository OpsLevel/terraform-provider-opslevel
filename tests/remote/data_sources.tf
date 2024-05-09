# Domain data sources

data "opslevel_domains" "all" {}

data "opslevel_domain" "first_domain_by_alias" {
  identifier = data.opslevel_domains.all.domains[0].aliases[0]
}

data "opslevel_domain" "first_domain_by_id" {
  identifier = data.opslevel_domains.all.domains[0].id
}

# Filter data sources

data "opslevel_filters" "all" {}

data "opslevel_filter" "first_filter_by_name" {
  filter {
    field = "name"
    value = data.opslevel_filters.all.filters[0].name
  }
}

data "opslevel_filter" "first_filter_by_id" {
  filter {
    field = "id"
    value = data.opslevel_filters.all.filters[0].id
  }
}

# Integration data sources

data "opslevel_integrations" "all" {}

data "opslevel_integration" "first_integration_by_id" {
  filter {
    field = "id"
    value = data.opslevel_integrations.all.integrations[0].id
  }
}

data "opslevel_integration" "first_integration_by_name" {
  filter {
    field = "name"
    value = data.opslevel_integrations.all.integrations[0].name
  }
}

# Lifecycle data sources

data "opslevel_lifecycles" "all" {}

data "opslevel_lifecycle" "first_lifecycle_by_alias" {
  filter {
    field = "alias"
    value = data.opslevel_lifecycles.all.lifecycles[0].alias
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_id" {
  filter {
    field = "id"
    value = data.opslevel_lifecycles.all.lifecycles[0].id
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_index" {
  filter {
    field = "index"
    value = data.opslevel_lifecycles.all.lifecycles[0].index
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_name" {
  filter {
    field = "name"
    value = data.opslevel_lifecycles.all.lifecycles[0].name
  }
}

# TODO: PropertyDefinition tests works on orange. Need to add to PAT acct.
# PropertyDefinition data sources

# data "opslevel_property_definitions" "all" {}

# data "opslevel_property_definition" "first_property_definition_by_id" {
#   identifier = data.opslevel_property_definitions.all.property_definitions[0].id
# }

# Rubric Category data sources

data "opslevel_rubric_categories" "all" {}

data "opslevel_rubric_category" "first_category_by_id" {
  filter {
    field = "id"
    value = data.opslevel_rubric_categories.all.rubric_categories[0].id
  }
}

data "opslevel_rubric_category" "first_category_by_name" {
  filter {
    field = "name"
    value = data.opslevel_rubric_categories.all.rubric_categories[0].name
  }
}

# Rubric Level data sources

data "opslevel_rubric_levels" "all" {}

data "opslevel_rubric_level" "first_level_by_alias" {
  filter {
    field = "alias"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].alias
  }
}

data "opslevel_rubric_level" "first_level_by_id" {
  filter {
    field = "id"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].id
  }
}

data "opslevel_rubric_level" "first_level_by_index" {
  filter {
    field = "index"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].index
  }
}

data "opslevel_rubric_level" "first_level_by_name" {
  filter {
    field = "name"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].name
  }
}

# Repository Category data sources

data "opslevel_repositories" "all" {}

data "opslevel_repository" "first_repo_by_alias" {
  alias = data.opslevel_repositories.all.repositories[0].alias
}

data "opslevel_repository" "first_repo_by_id" {
  id = data.opslevel_repositories.all.repositories[0].id
}

# Scorecard data sources

data "opslevel_scorecards" "all" {}

data "opslevel_scorecard" "first_scorecard_by_id" {
  identifier = data.opslevel_scorecards.all.scorecards[0].id
}

# Service data sources

data "opslevel_services" "all" {}

data "opslevel_service" "first_service_by_id" {
  id = data.opslevel_services.all.services[0].id
}

# System data sources

data "opslevel_systems" "all" {}

data "opslevel_system" "first_system_by_alias" {
  identifier = data.opslevel_systems.all.systems[0].aliases[0]
}

data "opslevel_system" "first_system_by_id" {
  identifier = data.opslevel_systems.all.systems[0].id
}

# Team data sources

data "opslevel_teams" "all" {}

data "opslevel_team" "first_team_by_alias" {
  alias = data.opslevel_teams.all.teams[0].alias
}

data "opslevel_team" "first_team_by_id" {
  id = data.opslevel_teams.all.teams[0].id
}

# Tier data sources

data "opslevel_tiers" "all" {}

# data "opslevel_tier" "first_tier_by_alias" {
#   filter {
#     field = "alias"
#     value = data.opslevel_tiers.all.tiers[0].alias
#   }
# }

data "opslevel_tier" "first_tier_by_id" {
  filter {
    field = "id"
    value = data.opslevel_tiers.all.tiers[0].id
  }
}

data "opslevel_tier" "first_tier_by_index" {
  filter {
    field = "index"
    value = data.opslevel_tiers.all.tiers[0].index
  }
}

data "opslevel_tier" "first_tier_by_name" {
  filter {
    field = "name"
    value = data.opslevel_tiers.all.tiers[0].name
  }
}

# User data sources

data "opslevel_users" "all" {}

data "opslevel_user" "first_user_by_email" {
  identifier = data.opslevel_users.all.users[0].email
}

data "opslevel_user" "first_user_by_id" {
  identifier = data.opslevel_users.all.users[0].id
}

# Webhook Action data sources

data "opslevel_webhook_actions" "all" {}

# data "opslevel_webhook_action" "first_webhook_action_by_alias" {
#   identifier = data.opslevel_webhook_actions.all.webhook_actions[0].aliases[0]
# }

data "opslevel_webhook_action" "first_webhook_action_by_id" {
  identifier = data.opslevel_webhook_actions.all.webhook_actions[0].id
}
