resource "opslevel_trigger_definition" "test" {
  access_control           = var.access_control
  action                   = var.action
  approval_required        = var.approval_required
  approval_users           = var.approval_users
  description              = var.description
  entity_type              = var.entity_type
  extended_team_access     = var.extended_team_access
  filter                   = var.filter
  manual_inputs_definition = var.manual_inputs_definition
  name                     = var.name
  owner                    = var.owner
  response_template        = var.response_template
  published                = var.published
}
