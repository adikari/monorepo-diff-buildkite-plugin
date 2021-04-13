FROM golang:1.16.3 as go

WORKDIR /plugin

COPY . .

RUN make clean build-linux

FROM buildkite/plugin-tester

COPY . .

COPY --from=go \
  /plugin/monorepo-diff-buildkite-plugin-linux \
  monorepo-diff-buildkite-plugin-linux
