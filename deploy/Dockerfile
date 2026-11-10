FROM golang:1.16 AS builder

WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-w" -o bin/manager \
    github.com/kubesphere/s2ioperator/cmd/manager

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /work/bin/manager manager
USER nonroot:nonroot
ENTRYPOINT ["/manager"]
