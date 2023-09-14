data "opslevel_rubric_category" "security" {
  filter {
    field = "name"
    value = "Security"
  }
}

data "opslevel_rubric_category" "Observability" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvNjA3"
  }
}