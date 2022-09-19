resource "opslevel_rubric_category" "example" {
  name = "foo"
}

output "category" {
  value = opslevel_rubric_category.example.id
}
