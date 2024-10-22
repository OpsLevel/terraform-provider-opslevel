resource "terraform_data" "opslevel_cli" {
  provisioner "local-exec" {
    command = "opslevel ${var.command}"
  }
}
