package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opslevel/opslevel-go/v2026"
)

func extractContactFromContacts(contactId opslevel.ID, contacts []opslevel.Contact) *opslevel.Contact {
	for _, readContact := range contacts {
		if contactId == readContact.Id {
			return &readContact
		}
	}
	return nil
}

func extractTagFromTags(tagId opslevel.ID, tags []opslevel.Tag) *opslevel.Tag {
	for _, readTag := range tags {
		if tagId == readTag.Id {
			return &readTag
		}
	}
	return nil
}

func extractToolFromTools(toolId opslevel.ID, tools []opslevel.Tool) *opslevel.Tool {
	for _, readTool := range tools {
		if toolId == readTool.Id {
			return &readTool
		}
	}
	return nil
}

func getTagsFromResource(client *opslevel.Client, resource opslevel.TaggableResourceInterface) (*opslevel.TagConnection, diag.Diagnostics) {
	var diags diag.Diagnostics
	tags, err := resource.GetTags(client, nil)
	if err != nil {
		diags.AddError(
			"opslevel client error",
			fmt.Sprintf("unable to read tags on %s with id '%s', got error: %s", string(resource.ResourceType()), string(resource.ResourceId()), err),
		)
	}
	if tags == nil {
		diags.AddError(
			"opslevel client error",
			fmt.Sprintf("zero tags found on %s with id '%s'", string(resource.ResourceType()), string(resource.ResourceId())),
		)
	}
	return tags, diags
}

// hasTagFormat returns true if the given tag is formatted as '<key>:<value>'
func hasTagFormat(tag string) bool {
	parts := strings.Split(tag, ":")
	return len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0
}

// isTagValid returns true if the given tag is formatted as '<resource-id>:<tag-id>'
func isTagValid(tag string) bool {
	ids := strings.Split(tag, ":")
	return len(ids) == 2 && opslevel.IsID(ids[0]) && opslevel.IsID(ids[1])
}
