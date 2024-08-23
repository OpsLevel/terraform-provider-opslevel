resource "opslevel_service_dependency" "test" {
  depends_upon = var.depends_upon
  note         = var.note
  service      = var.service
}
