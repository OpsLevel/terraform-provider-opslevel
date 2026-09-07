data "opslevel_teams" "all" {}

output "all" {
  value = data.opslevel_teams.all.teams
}

output "team_names" {
  value = sort(data.opslevel_teams.all.teams[*].name)
}

# Use OpsLevel teams as the source of truth for team metadata in other providers.
# Tags are returned as "key:value" strings, property values as JSON strings.
output "team_metadata" {
  value = {
    for team in data.opslevel_teams.all.teams : team.alias => {
      tags = team.tags
      properties = {
        for property in team.properties :
        property.definition.aliases[0] => jsondecode(property.value)
        if length(property.definition.aliases) > 0 && property.value != null
      }
    }
  }
}
