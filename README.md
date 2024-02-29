# Monorepo-diff-buildkite-plugin

This Monorepo plugin will assist you in triggering pipelines, as well as run commands in your CI by watching folders in your `monorepo`.

Check out this post to learn more about [**How to set up Continuous Integration for monorepo using Buildkite**](https://adikari.medium.com/set-up-continuous-integration-for-monorepo-using-buildkite-61539bb0ed76).

## Usefulness of the monorepo

A monorepo is a single, version-controlled code repository that houses multiple independent projects, offering benefits such as flexibility, streamlined management, and reduced tracking of changes and dependencies across multiple repositories.

This approach allows teams to:

- Reduce overhead associated with duplicating code for microservices.
- Easily maintain and monitor the entire codebase.

Check out the [example monorepo source code](https://github.com/buildkite/monorepo-example).

## Using the plugin

If the version number is not provided then the most recent version of the plugin will be used. Do not use version number as `master` or any branch names.

#### `watch`

It defines a list of paths or path to monitor for changes in the monorepo. It checks to see if there is a change to the subfolders specified in the path

#### `path`

A path or a list of paths to be watched, This part specifies which directory should be monitored. It can also be a glob pattern. For example specify `path: "**/*.md"` to match all markdown files. A list of paths can be provided to trigger the desired pipeline or run command or even do a pipeline upload.

#### `config`

This is a sub-section that provides configuration for running commands or triggering another pipeline when changes occur in the specified path
Configuration supports 2 different step types.

- [Trigger](https://buildkite.com/docs/pipelines/trigger-step)

  The configuration for the `trigger` step https://buildkite.com/docs/pipelines/trigger-step

  **Example**
  <br/>

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - monorepo-diff#v1.0.1:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: app/
              config:
                trigger: "app-deploy"
            - path: test/bin/
              config:
                command: "echo Make Changes to Bin"
```

- Changes to the path `app/` triggers the pipeline `app-deploy`
- Changes to the path `test/bin` will run the respective configuration command

 <br/>

⚠️ Warning : The user has to explictly state the paths they want to monitor. For instance if a user, is only watching path `app/` changes made to `app/bin` will not trigger the configuration. This is because the subfolder `/bin` was not specified.

 <br/>

**Example**
<br/>

```yaml
steps:
  - label: "Triggering pipelines with plugin"
    plugins:
      - monorepo-diff#v1.0.1:
          watch:
            - path: test/.buildkite/
              config: # Required [trigger step configuration]
                trigger: test-pipeline # Required [trigger pipeline slug]
            - path:
                - app/
                - app/bin/service/
              config:
                trigger: "data-generator"
                label: ":package: Generate data"
                build:
                  meta_data:
                    release-version: "1.1"
```

- When changes are detected in the path `test/.buildkite/` it triggers the pipeline `test-pipeline`
- If the changes are made to either `app/` or `app/bin/service/` it triggers the pipeline `data-generator`

<br/>

#### `diff` (optional)

This will run the script provided to determine the folder changes.
Depending on your use case, you may want to determine the point where the branch occurs
https://stackoverflow.com/questions/1527234/finding-a-branch-point-with-git and perform a diff against the branch point.

##### Sample output:

```
README.md
lib/trigger.bash
tests/trigger.bats
```

Default: `git diff --name-only HEAD~1`

##### Examples:

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

**Example**

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - monorepo-diff#v1.0.1:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: "bar-service/"
              config:
                command: "echo deploy-bar"
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
```

#### `interpolation` (optional)

This controls the pipeline interpolation on upload, and defaults to `true`.
If set to `false` it adds `--no-interpolation` to the `buildkite pipeline upload`,
to avoid trying to interpolate the commit message, which can cause failures.

#### `env` (optional)

The object values provided in this configuration will be appended to `env` property of all steps or commands.

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - monorepo-diff#v1.0.1:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
                label: "Triggered deploy"
                build:
                  message: "Deploying foo service"
                  env:
                    - HELLO=123
                    - AWS_REGION
```

#### `log_level` (optional)

Add `log_level` property to set the log level. Supported log levels are `debug` and `info`. Defaults to `info`.

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - monorepo-diff#v1.0.1:
          diff: "git diff --name-only HEAD~1"
          log_level: "debug" # defaults to "info"
          watch:
            - path: "foo-service/"
              config:
                trigger: "deploy-foo-service"
```

#### `hooks` (optional)

Currently supports a list of `commands` you wish to execute after the `watched` pipelines have been triggered

```yaml
hooks:
  - command: upload unit tests reports
  - command: echo success
```

#### `wait` (optional)

Default: `true`

By setting `wait` to `true`, the build will wait until the triggered pipeline builds are successful before proceeding

**Example**

```yaml
steps:
  - label: "Triggering pipelines"
    plugins:
      - monorepo-diff#v1.0.1:
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
                # following configs are available in command. notify is not available in trigger step
                notify:
                  - basecamp_campfire: https://basecamp-url
                  - github_commit_status:
                      context: my-custom-status
                  - slack: "@someuser"
                    if: build.state === "passed"
                # soft_fail: true
                soft_fail:
                  - exit_status: 1
                  - exit_status: "255"
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

          wait: true
```

## Thanks :heart:

Thanks to [@chronotc](https://github.com/chronotc) and [Monebag](https://github.com/monebag/) for authoring the original Buildkite Monorepo Plugin.

## License

MIT (see [LICENSE](LICENSE))
