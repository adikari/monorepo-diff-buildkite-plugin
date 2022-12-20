FROM golang:1.19-alpine as go

RUN apk add --no-cache git make \
  && go install github.com/goreleaser/goreleaser@latest

WORKDIR /plugin

COPY . .

RUN make build

FROM buildkite/plugin-tester

COPY . .

COPY --from=go \
  /plugin/dist/monorepo-diff-buildkite-plugin_linux_amd64_v1/monorepo-diff-buildkite-plugin \
  monorepo-diff-buildkite-plugin
