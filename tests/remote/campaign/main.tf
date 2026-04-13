resource "opslevel_campaign" "test" {
  name          = var.name
  owner_id      = var.owner_id
  project_brief = var.project_brief
  start_date    = var.start_date
  target_date   = var.target_date
}
