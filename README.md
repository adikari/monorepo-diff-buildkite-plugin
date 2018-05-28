# monorepo-diff-buildkite-plugin

This plugin will assist you in triggering pipelines by watching folders in your `monorepo`.

## Example

```yml
  steps:
    - label: "Triggering pipelines"
      plugins:
        monorepo-diff:
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
            - label: "Tag successful build"
              command: "echo "$(git rev-parse HEAD)" > last_successful_build"
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

git describe --tags --match production-* --abbrev=0
```

## References

https://stackoverflow.com/questions/1527234/finding-a-branch-point-with-git

## Contribute

### To run tests

Ensure that all tests are in the `./tests`

`docker-compose run --rm tests`