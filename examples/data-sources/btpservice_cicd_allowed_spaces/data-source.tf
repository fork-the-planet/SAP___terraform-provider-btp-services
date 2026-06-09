data "btpservice_cicd_allowed_spaces" "all" {}

output "allowed_space_guids" {
  value = [for s in data.btpservice_cicd_allowed_spaces.all.values : s.space_guid]
}
