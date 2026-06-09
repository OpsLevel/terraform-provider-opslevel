variable "name" {
  type = string
}

variable "owner_id" {
  type = string
}

variable "project_brief" {
  type    = string
  default = null
}

variable "start_date" {
  type    = string
  default = null
}

variable "target_date" {
  type    = string
  default = null
}

variable "reminder" {
  type = object({
    channels                        = list(string)
    frequency                       = number
    frequency_unit                  = string
    time_of_day                     = string
    timezone                        = string
    days_of_week                    = optional(list(string))
    message                         = optional(string)
    default_slack_channel           = optional(string)
    default_microsoft_teams_channel = optional(string)
  })
  default = null
}
