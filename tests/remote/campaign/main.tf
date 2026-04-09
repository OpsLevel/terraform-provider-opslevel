data "opslevel_campaigns" "all" {}

resource "opslevel_campaign" "test" {
  name          = var.name
  owner_id      = var.owner_id
  project_brief = var.project_brief
}
