# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.5.1]

### Changed
- Improve error logging by @adikari in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/118

## [2.5.0]

### Fixed
- Encode wait and hook steps with go-yaml by @nhurden in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/107
- Rebuild plugin_test in the Makefile by @nhurden in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/108

### Changed
- Update paths to new org url by @adikari in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/110

## [2.4.0]
- Support both "commands" and "command" key, of any type by @glittershark in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/104
- use log.Fatalf when logging error by @pmenglund in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/102
- Allow all tags when specifying the agent on a step by @uniquefine in https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/85

## [2.3.0]

### Added
- Added support for `soft_failure` step attribute. [#93](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/93) from [@isaacsu](https://github.com/isaacsu)

## [2.2.1]

### Fixed
- Detection of arm64 arch. [#84](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/84) from [@LColman-SEEK](https://github.com/LColman-SEEK)
- Panic when adjectent plugins used the same YAML fields. [#92](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/92) from [@isaacsu](https://github.com/isaacsu)

## [2.2.0]

### Added
- Added support for [Group Steps](https://buildkite.com/docs/pipelines/group-step). [#89](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/89) from [@xzyfer](https://github.com/xzyfer)

## [2.1.4]

### Fixed
- Fix failing MacOS builds. [#87](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/87) from [@xzyfer](https://github.com/xzyfer)

## [2.1.3]

### Fixed
- Log pipeline generation and upload errors. [#86](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/86) from [@xzyfer](https://github.com/xzyfer)

## [2.1.2]

### Fixed
- Prevent incorrect interpolation of env on trigger steps. [#80](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/80) from [@chronotc](https://github.com/monebag)

## [2.1.1]

### Added
- Support for linux-ppc64le platforms. [#78](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/78) from [@runlevel5](https://github.com/runlevel5)

## [2.1.0]

### Added
- Support for darwin-amd64 and darwin-arm64 platforms. [#74](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/74) from [@runlevel5](https://github.com/runlevel5)

## [2.0.5]

### Added
- Support for linux-arm64 platform. [#73](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/73) from [@runlevel5](https://github.com/runlevel5)
- Upgrade go version to 1.7.2 [#74](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/74) from [@runlevel5](https://github.com/runlevel5)

## [2.0.4]

### Added
- Glob support in path configuration. [#62](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/62) from [@worxli](https://github.com/worxli)

## [2.0.3]

### Fixed
- Get the appropriate version of binary based on what is specified in plugin configuration

## [2.0.2]

[#49](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/49) from [@adikari](https://github.com/adikari)

### Added
- Download the latest version of binary if plugin version is not specified
- Log plugin binary version

## [2.0.1]

### Fixed
- Incorrect binary path to download and docker service name [commit](https://github.com/monebag/monorepo-diff-buildkite-plugin/commit/a48753476822596c181d0f66cffb0d21fdc10214) from [@adikari](https://github.com/adikari)

## [2.0.0]

### Changed
- Rewrite plugin in go from [@adikari](https://github.com/adikari)

## [[1.3.2]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/1.3.2) - 2021-03-10

### Added
- Ability to pass `env` config in the `command` [#43](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/43) from [@jullianoacqio](https://github.com/jullianoacqio)

```yaml
...
   command: "echo $MESSAGE"
   env:
       MESSAGE: hello world
...
```

## [[1.3.1]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.3.1) - 2021-02-19

### Added
- Ability to pass `artificates` configuration [#41](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/41) from [@worxli](https://github.com/worxli)

```yaml
config:
   ...
   artifacts:
     - "logs/*"
```

### Fixed
- Fix parsing of queue name [#38](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/38) from [@ronaldmiranda](https://github.com/ronaldmiranda)

## [[1.3.0]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.3.0) - 2020-02-19

### Added
- Add support for leading emojis [#33](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/33) from [@worxli](https://github.com/worxli)

### Changed
- Extend command support with `label` and `queues` [#33](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/33) from [@worxli](https://github.com/worxli)

```yaml
steps:
  - label: "Deploying frontends"
    plugins:
      - monebag/monorepo-diff#v1.3.0:
          diff: "git diff --name-only HEAD~1"
          watch:
            - path: "react/"
              config:
                command: "netlify --production deploy"
                label: ":netlify: Deploy to production"
                agents:
                  queue: "deploy"
```

## [[1.2.0]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.2.0) - 2020-08-16

### Added
- Add support for `command` [#30](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/30) from [@chronotc](https://github.com/monebag)

```yaml
- monebag/monorepo-diff#v1.2.0:
          diff: "git diff --name-only $(head -n 1 last_successful_build)"
          watch:
            - path:
                - "ops/terraform/"
                - "ops/templates/terraform/"
              config:
                command: "buildkite-agent pipeline upload ops/.buildkite/pipeline.yml"
```


## [[1.1.1]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.1.1) - 2019-05-21

### Fixed
- Fix `has_changed` with large diff output [#13](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/13) from [@elliott-davis](https://github.com/elliott-davis)


## [[1.1.0]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.1.1) - 2019-02-19

### Added
- Ability to watch multiple paths and trigger a single pipeline [#10](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/10) from [@elliott-davis](https://github.com/elliott-davis)

### Changed
- Updated examples to be consistent with recommended Buildkite plugin syntax [#11](https://github.com/monebag/monorepo-diff-buildkite-plugin/pull/11) from [@harrietgrace](https://github.com/harrietgrace)


## [[1.0.0]](https://github.com/monebag/monorepo-diff-buildkite-plugin/releases/tag/v1.1.0) - 2018-10-04

### Added
- Initial Release from [@chronotc](https://github.com/monebag)
