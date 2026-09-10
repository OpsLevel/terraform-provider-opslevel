data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_team" "devs_by_id" {
  id = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NzQ0"
}

# Tags are returned as "key:value" strings.
output "devs_tags" {
  value = data.opslevel_team.devs.tags
}

# properties contains an entry for every team property definition, so filter to the
# ones this team actually has a value for. Values are JSON encoded.
output "devs_properties" {
  value = {
    for property in data.opslevel_team.devs.properties :
    property.definition.aliases[0] => jsondecode(property.value)
    if property.value != null && length(property.definition.aliases) > 0
  }
}
