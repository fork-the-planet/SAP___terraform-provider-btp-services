# Get the latest running (or most recently triggered) build.
data "btpservice_cicd_builds" "latest" {
  job    = "my-pipeline-job"
  filter = "latest"
}

# Get the most recently completed build.
data "btpservice_cicd_builds" "last_finished" {
  job    = "my-pipeline-job"
  filter = "latestFinished"
}

# Use the build ID to abort the currently running build.
resource "terraform_data" "abort_latest" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btpservice_cicd_abort_build.running]
    }
  }
}

action "btpservice_cicd_abort_build" "running" {
  config {
    job   = data.btpservice_cicd_builds.latest.job
    build = data.btpservice_cicd_builds.latest.builds[0].id
  }
}

# Use the build ID to delete the last finished build.
resource "terraform_data" "delete_last" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btpservice_cicd_delete_build.last]
    }
  }
}

action "btpservice_cicd_delete_build" "last" {
  config {
    job   = data.btpservice_cicd_builds.last_finished.job
    build = data.btpservice_cicd_builds.last_finished.builds[0].id
  }
}
