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

data "opslevel_integration" "kubernetes" {
    filter {
        field = "name"
        value = "Kubernetes"
    }
}

resource "opslevel_check_custom_event" "example" {
  name = "foo"
  enabled = true
  # To set a future enable date remove field 'enabled' and use 'enable_on'
  # enable_on = "2022-05-23T14:14:18.782000Z"
  category = data.opslevel_rubric_category.security.id
  level = data.opslevel_rubric_level.bronze.id
  owner = data.opslevel_team.devs.id
  filter = data.opslevel_filter.tier1.id
  integration = data.opslevel_integration.kubernetes.id
  service_selector = ".messages[] | .incident.service.id"
  success_condition = ".messages[] |   select(.incident.service.id == $ctx.alias) | .incident.status == \"resolved\""
  message = <<-EOT
  {% if check.passed %}
    ### Check passed
  {% else %}
    ### Check failed
    service **{{ data.messages[ctx.index].incident.service.id }}** has an unresolved incident.
  {% endif %}
  OpsLevel note: here you can fill in more details about this check. You can even include `data` from the payload, `params` specified in the URL and context `ctx` such as the service alias for the current evaluation.
  EOT
  notes = "Optional additional info on why this check is run or how to fix it"
}