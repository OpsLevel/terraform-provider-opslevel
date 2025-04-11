variable "access_control" {
  type        = string
  description = "The set of users that should be able to use the Trigger Definition."
}

variable "action" {
  type        = string
  description = "The action that will be triggered by the Trigger Definition."
}

variable "approval_required" {
  type    = bool
  default = false
}

variable "approval_teams" {
  type = list(string)
  default = []
}

variable "approval_users" {
  type = list(string)
  default = []
}

variable "description" {
  type        = string
  description = "The description of what the Trigger Definition will do."
  default     = null
}

variable "entity_type" {
  type        = string
  description = "The entity type to associate with the Trigger Definition."
  default     = null
}

variable "extended_team_access" {
  type        = list(string)
  description = "The set of additional teams who can invoke this Trigger Definition."
  default     = []
}

variable "filter" {
  type        = string
  description = "A filter defining which services this Trigger Definition applies to."
  default     = null
}

variable "manual_inputs_definition" {
  type        = string
  description = "The YAML definition of any custom inputs for this Trigger Definition."
  default     = null
}

variable "name" {
  type        = string
  description = "The name of the Trigger Definition"
}

variable "owner" {
  type        = string
  description = "The owner of the Trigger Definition."
}

variable "response_template" {
  type        = string
  description = "The liquid template used to parse the response from the Webhook Action."
  default     = null
}

variable "published" {
  type        = bool
  description = "The published state of the Custom Action; true if the Trigger Definition is ready for use; false if it is a draft."
}
