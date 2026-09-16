# Import using the id of the check on the campaign - not the id of the rubric
# check it was copied from. The campaign is derived from the check itself.
#
# Find the ids of a campaign's checks with the opslevel_campaign data source:
#
#   data "opslevel_campaign" "example" {
#     identifier = "Z2lkOi8vb3BzbGV2ZWwvQ2FtcGFpZ24vMTIz"
#   }
#
#   output "campaign_checks" {
#     value = data.opslevel_campaign.example.checks
#   }
#
# copied_from_check_id is left null on an imported check: OpsLevel has no record
# of where copies made before it began tracking lineage came from, and the value
# cannot be set after the fact.
terraform import opslevel_campaign_check.example Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpUb29sVXNhZ2UvNDU2
