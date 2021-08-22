data "opslevel_lifecycle" "beta" {
    filter {
        field = "alias"
        value = "beta"
    }
}

data "opslevel_lifecycle" "pre-alpha" {
    filter {
        field = "index"
        value = "0"
    }
}