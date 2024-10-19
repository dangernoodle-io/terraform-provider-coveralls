resource "coveralls_repository" "example" {
  name                                = "dangernoodle-io/terraform-provider-coveralls"
  service                             = "github"
  comment_on_pull_requests            = true
  send_build_status                   = true
  commit_status_fail_threshold        = 3.7
  commit_status_fail_change_threshold = 5.0
}
