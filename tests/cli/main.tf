resource "terraform_data" "opslevel_cli" {
  provisioner "local-exec" {
    command = var.command == null ? ":" : "opslevel ${var.command}"
  }
}

resource "terraform_data" "delete_command" {
  provisioner "local-exec" {
    # NOTE: ':' is a no-op bash command
    command = var.resource_type == null ? ":" : "opslevel delete ${var.resource_type} ${var.resource_id}"
  }
}
