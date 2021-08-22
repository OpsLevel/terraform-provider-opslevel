# terraform-provider-opslevel

Terraform Provider for [OpsLevel](https://opslevel.com)



# Useful Snippets

Get Tag Values

```hcl
values({for entry in data.opslevel_service.example.tags : entry => split(":", entry)[1]})
```