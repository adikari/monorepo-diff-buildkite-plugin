FROM goreleaser/goreleaser:v1.13.1-amd64 as go

WORKDIR /plugin

COPY . .

RUN goreleaser build --rm-dist --skip-validate --config .goreleaser-test.yml

FROM buildkite/plugin-tester:v3.0.1

COPY tests /plugin/tests
COPY hooks /plugin/hooks

COPY --from=go \
  /plugin/dist/monorepo-diff-buildkite-plugin_linux_amd64_v1/monorepo-diff-buildkite-plugin \
  monorepo-diff-buildkite-plugin
