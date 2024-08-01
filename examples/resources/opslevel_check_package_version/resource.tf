data "opslevel_rubric_category" "security" {
  filter {
    field = "name"
    value = "Security"
  }
}

data "opslevel_rubric_level" "bronze" {
  filter {
    field = "name"
    value = "Bronze"
  }
}

data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_filter" "tier1" {
  filter {
    field = "name"
    value = "Tier 1"
  }
}

resource "opslevel_check_package_version" "example" {
  name     = "foo"
  enable   = true
  category = data.opslevel_rubric_category.security.id
  level    = data.opslevel_rubric_level.bronze.id
  notes    = "Optional additional info on why this check is run or how to fix it"

  package_constraint = "exists"
  package_manager    = "gradle"
  package_name       = "log4j"
}

resource "opslevel_check_package_version" "example2" {
  name     = "foo"
  enable   = true
  category = data.opslevel_rubric_category.security.id
  level    = data.opslevel_rubric_level.bronze.id
  notes    = "Optional additional info on why this check is run or how to fix it"

  package_constraint     = "matches-version"
  package_manager        = "npm"
  package_name           = "leftpad"
  missing_package_result = "passed"
  version_constraint_predicate = {
    type  = "matches_regex"
    value = "1.0.*"
  }
}

resource "opslevel_check_package_version" "example3" {
  name     = "foo"
  enable   = true
  category = data.opslevel_rubric_category.security.id
  level    = data.opslevel_rubric_level.bronze.id
  notes    = "Optional additional info on why this check is run or how to fix it"

  package_constraint     = "matches-version"
  package_manager        = "go"
  package_name           = "client-go/.*"
  package_name_is_regex  = true
  missing_package_result = "passed"
  version_constraint_predicate = {
    type  = "matches_regex"
    value = "1.27.*"
  }
}
