resource "opslevel_check_manual" "test" {
  category                = var.category
  enable_on               = var.enable_on
  filter                  = var.filter
  level                   = var.level
  name                    = var.name
  notes                   = var.notes
  owner                   = var.owner
  update_frequency        = var.update_frequency
  update_requires_comment = var.update_requires_comment

  lifecycle {
    ignore_changes = [
      enabled,
      enable_on
    ]
  }
}
