data "opslevel_team" "devs" {
  alias = "developers"
}

data "opslevel_team" "devs_by_id" {
  id = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS83NzQ0"
}

# Tags are returned as "key:value" strings, property values as JSON strings.
output "devs_tags" {
  value = data.opslevel_team.devs.tags
}

output "devs_properties" {
  value = {
    for property in data.opslevel_team.devs.properties :
    property.definition.aliases[0] => jsondecode(property.value)
    if length(property.definition.aliases) > 0 && property.value != null
  }
}