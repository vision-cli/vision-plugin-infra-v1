variable "environment" {
  description = "The environment this project is being created in, either dev or prod"
  type        = string
}

variable "org_id" {
  description = "The organization id for the associated services"
  type        = string
}

variable "folder_id" {
  description = "The folder where the projects will be created"
  type        = string
}

variable "project_name" {
  description = "The name of this project"
  type        = string
}

variable "unique_str" {
  description = "A unique identifier for this project"
  type        = string
}

variable "billing_account" {
  description = "The ID of the billing account to associate this project with"
  type        = string
}

variable "region" {
  description = "default location for resources"
  type        = string
  default     = "europe-west2"
}

variable "zone" {
  description = "default zone for resources in the default location"
  type        = string
  default     = "europe-west2-c"
}

variable "oauth2_client_id" {
  description = "Oauth2 client id"
  type        = string
}

variable "oauth2_client_secret" {
  description = "Oauth2 client secret"
  type        = string
}

variable "domain" {
  description = "Domain name to run the load balancer on. Used because `ssl` is `true`."
  type        = string
}

variable "members" {
  description = "List of members and groups who have access"
}

variable "db_user_name" {
  description = "The default username for the database"
  type        = string
}

variable "db_user_password" {
  description = "The default password for the database"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Default database name to create"
  type        = string
}

variable "db_tier" {
  type    = string
  default = "db-g1-small"
}

variable "db_vol_size" {
  type    = number
  default = 20
}

