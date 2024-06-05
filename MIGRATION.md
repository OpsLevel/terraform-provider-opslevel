# State Migration (v0.11.0 -> v1.0.0)

The structure of some resource data has changed from v0.11.0 to v1.0.0.
If you encounter erros like `Error: Unable to Read Previously Saved State for UpgradeResourceState`
while upgrading to the OpsLevel v1.0.0 Terraform provider, the following will get you unstuck.

### Get ids of resources that raise errors

Assuming `opslevel_check_tag_defined.my_tag` raises an error like
`AttributeName("environment_predicate"): invalid JSON, expected "{", got "["`, run
`terraform state show opslevel_check_tag_defined.my_tag` and extract the value of `id`.
Extract and keep all ids for resources that raise errors.

If `terraform state show ...` raises an error, try `terraform state pull > my_tf_state.json`
and extract ids from `my_tf_state.json`.

### Remove problematic resources from state

This does NOT remove the resources from OpsLevel - only from Terraform's state.
Run `terraform state rm opslevel_check_tag_defined.my_tag` and repeat for all
resources that raise errors.
Optionally verify resources removed from state with `terraform state list`.

### Import resources from OpsLevel into Terraform state

Import OpsLevel resources into Terraform's state by running
`terraform import opslevel_check_tag_defined.my_tag <resource-id>` and repeat
for all resources that raise errors. Note <resource-id> matches the `id` of the OpsLevel
resource as returned by the [OpsLevel cli](https://github.com/opslevel/cli), if `id` is
not otherwise known.
Optionally verify resources imported to state with `terraform state list`.
