# monorepo-diff-buildkite-plugin [![Build status](https://badge.buildkite.com/719d0b895285367c9c57a09e07f1e853148d2509f0667e0ae8.svg?branch=master)](https://buildkite.com/kuda/monorepo-diff-buildkite-plugin)

This plugin will assist you in triggering pipelines by watching folders in your `monorepo`.

## Example

### Simple

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - chronotc/monorepo-diff#v1.3.1:
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
      - chronotc/monorepo-diff#v1.3.1:
          diff: "git diff --name-only $(head -n 1 last_successful_build)"
          interpolation: false
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
          hooks:
            - command: "echo $(git rev-parse HEAD) > last_successful_build"
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

### `interpolation` (optional)

This controls the pipeline interpolation on upload, and defaults to `true`.
If set to `false` it adds `--no-interpolation` to the `buildkite pipeline upload`,
to avoid trying to interpolate the commit message, which can cause failures.

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
- path: app/cms/
  config:
    command: "netlify --production deploy"
    label: ":netlify: Deploy to production"
    agents:
      queue: "deploy"
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

## Environment

### `DEBUG` (optional)

By turning `DEBUG` on, the generated pipeline will be displayed prior to upload

```yaml
steps:
  - label: "Triggering pipelines"
    env:
      DEBUG: true
    plugins:
      - chronotc/monorepo-diff:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
```

## References

https://stackoverflow.com/questions/1527234/finding-a-branch-point-with-git

## Contribute

### To run tests

Ensure that all tests are in the `./tests`

`docker-compose run --rm tests`

### To run lint

`docker-compose run --rm lint`
