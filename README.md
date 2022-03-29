[![e2e status](https://badge.buildkite.com/719d0b895285367c9c57a09e07f1e853148d2509f0667e0ae8.svg?branch=master)](https://buildkite.com/kuda/monorepo-diff-buildkite-plugin)
[![codecov](https://codecov.io/gh/chronotc/monorepo-diff-buildkite-plugin/branch/master/graph/badge.svg?token=DQ3B4FIYD2)](https://codecov.io/gh/chronotc/monorepo-diff-buildkite-plugin)
[![Publish](https://github.com/chronotc/monorepo-diff-buildkite-plugin/actions/workflows/publish.yml/badge.svg)](https://github.com/chronotc/monorepo-diff-buildkite-plugin/actions/workflows/publish.yml)<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-11-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
# monorepo-diff-buildkite-plugin

This plugin will assist you in triggering pipelines by watching folders in your `monorepo`.

Check out this post to learn [**How to set up Continuous Integration for monorepo using Buildkite**](https://adikari.medium.com/set-up-continuous-integration-for-monorepo-using-buildkite-61539bb0ed76).

## Using the plugin

If the version number is not provided then the most recent version of the plugin will be used. Do not use version number as `master` or any branch names.

### Simple

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - chronotc/monorepo-diff#v2.1.4:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: "bar-service/"
              config:
                command: "echo deploy-bar"
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
```

### Detailed

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - chronotc/monorepo-diff#v2.1.4:
          diff: "git diff --name-only $(head -n 1 last_successful_build)"
          interpolation: false
          env:
            - env1=env-1 # this will be appended to all env configuration
          hooks:
            - command: "echo $(git rev-parse HEAD) > last_successful_build"
          watch:
            - path:
                - "ops/terraform/"
                - "ops/templates/terraform/"
              config:
                command: "buildkite-agent pipeline upload ops/.buildkite/pipeline.yml"
                label: "Upload pipeline"
                retry:
                  automatic:
                  - limit: 2
                    exit_status: -1
                agents:
                  queue: performance
                artifacts:
                  - "logs/*"
                env:
                  - FOO=bar
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
                label: "Triggered deploy"
                build:
                  message: "Deploying foo service"
                  env:
                    - HELLO=123
                    - AWS_REGION

          wait: true
```

## Configuration

## `diff` (optional)

This will run the script provided to determine the folder changes.
Depending on your use case, you may want to determine the point where the branch occurs
https://stackoverflow.com/questions/1527234/finding-a-branch-point-with-git and perform a diff against the branch point.

#### Sample output:
```
README.md
lib/trigger.bash
tests/trigger.bats
```

Default: `git diff --name-only HEAD~1`

#### Examples:

`diff: ./diff-against-last-successful-build.sh`

```bash
#!/bin/bash

set -ueo pipefail

LAST_SUCCESSFUL_BUILD_COMMIT="$(aws s3 cp "${S3_LAST_SUCCESSFUL_BUILD_COMMIT_PATH}" - | head -n 1)"
git diff --name-only "$LAST_SUCCESSFUL_BUILD_COMMIT"
```

`diff: ./diff-against-last-built-tag.sh`

```bash
#!/bin/bash

set -ueo pipefail

LATEST_BUILT_TAG=$(git describe --tags --match foo-service-* --abbrev=0)
git diff --name-only "$LATEST_TAG"
```

## `interpolation` (optional)

This controls the pipeline interpolation on upload, and defaults to `true`.
If set to `false` it adds `--no-interpolation` to the `buildkite pipeline upload`,
to avoid trying to interpolate the commit message, which can cause failures.

## `env` (optional)

The object values provided in this configuration will be appended to `env` property of all steps or commands.

## `log_level` (optional)

Add `log_level` property to set the log level. Supported log levels are `debug` and `info`. Defaults to `info`.

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - chronotc/monorepo-diff#v2.1.4:
          diff: "git diff --name-only HEAD~1"
          log_level: "debug" # defaults to "info"
          watch:
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
```

## `watch`

Declare a list of

```yaml
- path: app/cms/
  config: # Required [trigger step configuration]
    trigger: cms-deploy # Required [trigger pipeline slug]
- path:
    - services/email
    - assets/images/email
  config:
    trigger: email-deploy
```

### `path`

If the `path` specified here in the appears in the `diff` output, a `trigger` step will be added to the dynamically generated pipeline.yaml

A list of paths can be provided to trigger the desired pipeline. Changes in any of the paths will initiate the pipeline provided in trigger.

A `path` can also be a glob pattern. For example specify `path: "**/*.md"` to match all markdown files.

### `config`

Configuration supports 2 different step types.

- [Trigger](https://buildkite.com/docs/pipelines/trigger-step)
- [Command](https://buildkite.com/docs/pipelines/command-step)

#### Trigger

The configuration for the `trigger` step https://buildkite.com/docs/pipelines/trigger-step

By default, it will pass the following values to the `build` attributes unless an alternative values are provided

```yaml
- path: app/cms/
  config:
    trigger: cms-deploy
    build:
      commit: $BUILDKITE_COMMIT
      branch: $BUILDKITE_BRANCH
```

### `wait` (optional)

Default: `true`

By setting `wait` to `true`, the build will wait until the triggered pipeline builds are successful before proceeding

### `hooks` (optional)

Currently supports a list of `commands` you wish to execute after the `watched` pipelines have been triggered

```yaml
hooks:
  - command: upload unit tests reports
  - command: echo success
```

#### Command

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - chronotc/monorepo-diff#v2.1.4:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: app/cms/
              config:
                group: ":rocket: deployment"
                command: "netlify --production deploy"
                label: ":netlify: Deploy to production"
                agents:
                  queue: "deploy"
                env:
                  - FOO=bar
```

There is currently limited support for command configuration. Only the `command` property can be provided at this point in time.

Using commands, it is also possible to use this to upload other pipeline definitions

```yaml
- path: frontend/
  config:
    command: "buildkite-agent pipeline upload ./frontend/.buildkite/pipeline.yaml"
- path: infrastructure/
  config:
    command: "buildkite-agent pipeline upload ./infrastructure/.buildkite/pipeline.yaml"
- path: backend/
  config:
    command: "buildkite-agent pipeline upload ./backend/.buildkite/pipeline.yaml"
```

## How to Contribute

Please read [contributing guide](https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/master/CONTRIBUTING.md).

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://www.subash.com.au"><img src="https://avatars.githubusercontent.com/u/1757714?v=4?s=100" width="100px;" alt=""/><br /><sub><b>subash adhikari</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=adikari" title="Code">💻</a> <a href="#example-adikari" title="Examples">💡</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=adikari" title="Documentation">📖</a> <a href="#maintenance-adikari" title="Maintenance">🚧</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/pulls?q=is%3Apr+reviewed-by%3Aadikari" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=adikari" title="Tests">⚠️</a> <a href="#infra-adikari" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="https://github.com/chronotc"><img src="https://avatars.githubusercontent.com/u/7519144?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Silla Tan</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=chronotc" title="Code">💻</a> <a href="#example-chronotc" title="Examples">💡</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=chronotc" title="Documentation">📖</a> <a href="#maintenance-chronotc" title="Maintenance">🚧</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/pulls?q=is%3Apr+reviewed-by%3Achronotc" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=chronotc" title="Tests">⚠️</a> <a href="#infra-chronotc" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="http://excellent.io"><img src="https://avatars.githubusercontent.com/u/1130349?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Elliott Davis</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=elliott-davis" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=elliott-davis" title="Tests">⚠️</a> <a href="#ideas-elliott-davis" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="http://www.acqio.com.br"><img src="https://avatars.githubusercontent.com/u/35783178?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Julliano Gonçalves</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=jullianoacqio" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=jullianoacqio" title="Tests">⚠️</a></td>
    <td align="center"><a href="http://worx.li"><img src="https://avatars.githubusercontent.com/u/2504856?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Lukas Bischofberger</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=worxli" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=worxli" title="Tests">⚠️</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=worxli" title="Documentation">📖</a> <a href="#example-worxli" title="Examples">💡</a></td>
    <td align="center"><a href="http://blog.englund.nu/"><img src="https://avatars.githubusercontent.com/u/32618?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Martin Englund</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=pmenglund" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=pmenglund" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/jacekszubert"><img src="https://avatars.githubusercontent.com/u/17125006?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Jacek Szubert</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=jacekszubert" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="http://ronaldmiranda.com.br"><img src="https://avatars.githubusercontent.com/u/14340100?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Ronald Carvalho</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/issues?q=author%3Aronaldmiranda" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/harrietgrace"><img src="https://avatars.githubusercontent.com/u/2074469?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Harriet Lawrence</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=harrietgrace" title="Documentation">📖</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=harrietgrace" title="Examples">💡</a></td>
    <td align="center"><a href="https://github.com/runlevel5"><img src="https://avatars.githubusercontent.com/u/135605?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Trung Lê</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=runlevel5" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=runlevel5" title="Examples">💡</a><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=runlevel5" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/jquick"><img src="https://avatars.githubusercontent.com/u/574637?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Jared Quick</b></sub></a><br /><a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/issues?q=author%3jquick" title="Bug reports">🐛</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=jquick" title="Code">💻</a> <a href="https://github.com/chronotc/monorepo-diff-buildkite-plugin/commits?author=jquick" title="Tests">⚠️</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->
