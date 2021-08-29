module github.com/opslevel/terraform-provider-opslevel

go 1.16

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.4.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/kr/pretty v0.2.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/opslevel/opslevel-go v0.3.0
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// Uncomment for local development
// replace github.com/opslevel/opslevel-go => ../../go/src/github.com/opslevel/opslevel-go/
