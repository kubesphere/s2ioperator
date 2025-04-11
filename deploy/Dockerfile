FROM golang:1.16 AS builder

RUN apt update -y && apt install -y ca-certificates libgnutls30

WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-w" -o bin/manager \
    github.com/kubesphere/s2ioperator/cmd/manager

FROM alpine:3.21

WORKDIR /
COPY --from=builder /work/bin/manager manager
ENTRYPOINT ["/manager"]
