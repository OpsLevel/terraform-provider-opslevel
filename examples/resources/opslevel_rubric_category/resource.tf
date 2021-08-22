resource "opslevel_rubric_category" "example" {
  name = "foo"
  description = "foo category"
}

output "category" {
  value = opslevel_rubric_category.example.id
}
