resource "terraform_data" "delete-resource" {
  provisioner "local-exec" {
    command = "opslevel delete ${var.resource_type} ${var.resource_id}"
  }
}
