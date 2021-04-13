FROM golang:1.16

WORKDIR /go/src/github.com/chronotc/monorepo-diff-buildkite-plugin/

COPY . .

RUN make clean
RUN make build-linux

FROM buildkite/plugin-tester

COPY . .

COPY --from=builder /go/src/github.com/jwplayer/buildpipe-buildkite-plugin/buildpipe-linux buildpipe

CMD ["bats", "tests/"]
