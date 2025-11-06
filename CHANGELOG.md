# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-11-06
### Added
- **New Resource**: `slack_usergroup` - Create and manage Slack usergroups (user groups)
  - Support for creating, updating, and disabling usergroups
  - Manage usergroup members declaratively
  - Set handle, name, and description
- **New Feature**: `strict_members` parameter for `slack_channel` resource
  - When `false` (default): Terraform only manages declared members, ignores manually added users
  - When `true`: Terraform detects drift when users are manually added to channels
- Helper function library for improved code reusability

### Fixed
- **Critical Bug Fix**: Fixed panic when creating channels with members (#9)
  - Resolved type assertion error with `schema.Set` handling
  - Added proper type conversion helpers
  - Applied fix to both `resource_channel_create.go` and `resource_channel_update.go`

### Changed
- Improved member management in channels with better drift detection options
- Enhanced documentation with more examples and use cases
- Added `usergroups:write` to required OAuth scopes documentation

### Documentation
- Added comprehensive documentation for `slack_usergroup` resource
- Added examples for usergroup management
- Updated provider documentation with new features
- Clarified member management behavior and limitations

## [0.1.7] - 2025-XX-XX
### Changed
- Previous release changes

## [0.1.0] - 2025-06-15
### Added
- Initial release with support for:
  - Slack channels (create, update, purpose, topic, members)
  - Slack users and user lookup
  - Slack usergroups (data source)
  - Slack channels data source
