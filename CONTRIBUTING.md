# Contributing

First of all, thank you for your interest in contributing to this project.

Before creating a pull request, please read and follow this contributing guide.
Before participating in this project's activities, please read and follow [code of conduct](https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CODE_OF_CONDUCT.md).

Please create an [issue](https://github.com/chronotc/monorepo-diff-buildkite-plugin/issues) before submitting a pull request. If it is a proposal for a new feature or changing existing functionality, initiate a discussion with maintainers first. If it's a fix for know bugs, a discussion is not required.

## Developing

- [Install Go.](https://golang.org/doc/install)
- The project uses Makefile. Install `make` command.
- Fork this repository.
- Clone the forked repository.
-  Make changes (see [Formatting](https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CONTRIBUTING.md#formatting)) and commit to your fork. Commit messages should follow the [Conventional Commits](https://www.conventionalcommits.org/) style.
- Add appropriate unit tests (see [Testing](https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CONTRIBUTING.md#testing)) for your changes.
- Update documentation if appropriate.
- Create a pull request with your changes.
- Github action will run the necessary checks against your pull request.
- A maintainer will review the pull request once all checks are
- A maintainer will merge and create a release (see [Releasing](https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CONTRIBUTING.md#releasing)).

## Testing

All changes must be unit tested and meet the project test coverage threshold (73%) requirement.
Run `make test` to run all tests and generate coverage reports before submitting a pull request.

To write the `bats` tests for plugin,
1. Modify the tests
2. Run `docker-compose build plugin_test && docker-compose run --rm plugin_test`

## Formatting

All code must be formatted with `gofmt` (with the latest Go version) and pass `go vet`. The plugin must be linted with [buildkite-plugin-linter](https://github.com/buildkite-plugins/buildkite-plugin-linter). Run `make quality` to run all formatting checks.

## Releasing

One of the maintainers will create a release after merging the pull request.
- Ensure documentation is updated appropriately.
- Update [README.md]( https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/README.md ) and [hooks/command]( https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/hooks/command ) to the new version number.
- Update [ changelog ]( https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CHANGELOG.md ) using guidelines from [Keep a Changelog](https://keepachangelog.com/).
- Merge changes to master.
- Create a new release. The version number must follow the format of `v<major>.<minor>.<patch>`. Eg. `v2.0.1`
- Github actions automatically builds and attachs the binary to the release.
