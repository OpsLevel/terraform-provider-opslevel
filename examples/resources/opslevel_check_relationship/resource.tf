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

# Example: Create a relationship check that ensures services have at least 2 "depends_on" relationships
resource "opslevel_check_relationship" "example" {
  name        = "Service Dependencies Check"
  enabled     = true
  # To set a future enable date remove field 'enabled' and use 'enable_on'
  # enable_on = "2022-05-23T14:14:18.782000Z"
  category    = data.opslevel_rubric_category.security.id
  level       = data.opslevel_rubric_level.bronze.id
  owner       = data.opslevel_team.devs.id
  filter      = data.opslevel_filter.tier1.id
  notes       = "Ensures services have proper dependency relationships defined"
  
  # The relationship definition ID for "depends_on" relationships
  relationship_definition_id = "Z2lkOi8vc2VydmljZS8xMjM0NTY3ODk"  # Replace with actual relationship definition ID
  
  # Predicate to check that the count is greater than or equal to 2
  relationship_count_predicate {
    type  = "greater_than_or_equal_to"
    value = "2"
  }
}

# Example: Create a relationship check that ensures services have no more than 5 "belongs_to" relationships
resource "opslevel_check_relationship" "example_max_relationships" {
  name        = "Service Belongs To Check"
  enabled     = true
  category    = data.opslevel_rubric_category.security.id
  level       = data.opslevel_rubric_level.bronze.id
  notes       = "Ensures services don't have too many 'belongs_to' relationships"
  
  # The relationship definition ID for "belongs_to" relationships
  relationship_definition_id = "Z2lkOi8vc2VydmljZS8xMjM0NTY3ODk"  # Replace with actual relationship definition ID
  
  # Predicate to check that the count is less than or equal to 5
  relationship_count_predicate {
    type  = "less_than_or_equal_to"
    value = "5"
  }
}
