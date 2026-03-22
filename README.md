# terraform-provider-coveralls
[![provider-test](https://github.com/dangernoodle-io/terraform-provider-coveralls/actions/workflows/provider-test.yml/badge.svg?branch=main)](https://github.com/dangernoodle-io/terraform-provider-coveralls/actions/workflows/provider-test.yml)
[![provider-release](https://github.com/dangernoodle-io/terraform-provider-coveralls/actions/workflows/provider-release.yml/badge.svg)](https://github.com/dangernoodle-io/terraform-provider-coveralls/actions/workflows/provider-release.yml)
[![GitHub Release](https://img.shields.io/github/v/release/dangernoodle-io/terraform-provider-coveralls)](https://github.com/dangernoodle-io/terraform-provider-coveralls/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/dangernoodle-io/terraform-provider-coveralls/badge.svg?branch=main)](https://coveralls.io/github/dangernoodle-io/terraform-provider-coveralls?branch=main)
[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-dangernoodle--io%2Fcoveralls-7B42BC?logo=terraform)](https://registry.terraform.io/providers/dangernoodle-io/coveralls/latest/docs)
[![Go version](https://img.shields.io/github/go-mod/go-version/dangernoodle-io/terraform-provider-coveralls)](go.mod)

A Terraform provider for managing [Coveralls](https://coveralls.io) repositories.

## Resources

### `coveralls_repository`

Manages a Coveralls repository configuration.

```terraform
resource "coveralls_repository" "example" {
  name                                = "dangernoodle-io/terraform-provider-coveralls"
  service                             = "github"
  comment_on_pull_requests            = true
  send_build_status                   = true
  commit_status_fail_threshold        = 3.7
  commit_status_fail_change_threshold = 5.0
}
```

#### Arguments

- `name` - (Required) Repository name in `owner/repo` format.
- `service` - (Required) Source control service (e.g. `github`).
- `comment_on_pull_requests` - (Required) Whether to post comments on pull requests.
- `send_build_status` - (Required) Whether to send build status to the source control service.
- `commit_status_fail_threshold` - (Optional) Coverage threshold below which to fail the build.
- `commit_status_fail_change_threshold` - (Optional) Coverage change threshold below which to fail the build.

#### Attributes

- `token` - Coveralls repository token.
- `created_at` - Timestamp of when the repository was created.
- `updated_at` - Timestamp of when the repository was last updated.

#### Import

```shell
terraform import coveralls_repository.example github:dangernoodle-io/terraform-provider-coveralls
```

## Data Sources

### `coveralls_repository`

Reads an existing Coveralls repository.

```terraform
data "coveralls_repository" "example" {
  name    = "dangernoodle-io/terraform-provider-coveralls"
  service = "github"
}
```

#### Arguments

- `name` - (Required) Repository name in `owner/repo` format.
- `service` - (Required) Source control service (e.g. `github`).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (for building)

## Developing

```shell
make          # fmt, lint, install, generate
make test     # unit tests
make testacc  # acceptance tests (requires TF_ACC=1 and COVERALLS_API_TOKEN)
```