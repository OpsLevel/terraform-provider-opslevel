data "opslevel_rubric_level" "bronze" {
  filter {
    field = "alias"
    value = "bronze"
  }
}

data "opslevel_rubric_level" "gold" {
  filter {
    field = "id"
    value = "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvMzE5"
  }
}
