# monorepo-diff-buildkite-plugin

This plugin will assist you in triggering pipelines by watching folders in your `monorepo`.

## Example

### Simple

```yml
steps:
  - label: "Triggering pipelines"
    plugins:
      chronotc/monorepo-diff#v1.0.0:
        diff: "git diff --name-only HEAD~1"
        watch:
          - path: "foo-service/"
            config:
              trigger: "deploy-foo-service"
```

### Detailed

```yml
steps:
  - label: "Triggering pipelines"
    plugins:
      chronotc/monorepo-diff#v1.0.0:
        diff: "git diff --name-only $(head -n 1 last_successful_build)"
        watch:
          - path: "foo-service/"
            config:
              trigger: "deploy-foo-service"
              build:
                message: "Deploying foo service"
                env:
                  - HELLO=123
                  - AWS_REGION
          - path: "ops/terraform/"
            config:
              trigger: "provision-terraform-resources"
              async: true
        wait: true
        hooks:
          - command: "echo $(git rev-parse HEAD) > last_successful_build"
```

## Configuration

### `diff` (optional)

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

### `watch`

Declare a list of

```yaml
- path: app/cms/
  config: # Required [trigger step configuration]
    trigger: cms-deploy # Required [trigger pipeline slug]
```

#### `path`

If the `path` specified here in the appears in the `diff` output, a `trigger` step will be added to the dynamically generated pipeline.yml

#### `config`

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

By setting `wait` to `true`, the build will wait until the triggered pipeline builds are successful before proceeding

### `hooks` (optional)

Currently supports a list of `commands` you wish to execute after the `watched` pipelines have been triggered

```yaml
hooks:
  - command: upload unit tests reports
  - command: echo success

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
      chronotc/monorepo-diff:
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