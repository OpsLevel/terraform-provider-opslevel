package opslevel

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var opslevelPropertyObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"definition": types.StringType,
		"owner":      types.StringType,
		"value":      types.StringType,
	},
}
