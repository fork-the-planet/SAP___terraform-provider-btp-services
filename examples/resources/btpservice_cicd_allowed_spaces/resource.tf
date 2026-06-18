resource "btpservice_cicd_allowed_spaces" "this" {
  allowed_spaces = [
    {
      space_guid = "a2bcf2b8-6eda-5b8a-0b7c-8512bb82060f"
      comment    = "Team Alpha production space"
    },
    {
      space_guid = "d794d687-3053-4cba-a942-88e6b13ef035"
      comment    = "Team Beta staging space"
    },
  ]
}
