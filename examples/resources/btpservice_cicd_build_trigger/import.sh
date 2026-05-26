# terraform import btpservice_cicd_build_trigger.<resource_name> <job_id>,<trigger_id>

terraform import btpservice_cicd_build_trigger.example my-pipeline-job,afbacb1c-e7f2-4f8d-94e3-8508332fcd2e

# terraform import using id attribute in import block

import {
  to = btpservice_cicd_build_trigger.<resource_name>
  id = "<job_id>,<trigger_id>"
}

import {
  to = btpservice_cicd_build_trigger.<resource_name>
  identity = {
    job = "<job_id>"
    id  = "<trigger_id>"
  }
}
