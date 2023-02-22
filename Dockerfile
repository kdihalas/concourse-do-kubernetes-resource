FROM golang:1.20.1-bullseye AS builder

WORKDIR /concourse/concourse-resource
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV CGO_ENABLED 0
RUN go build -o /assets/check github.com/kdihalas/concourse-do-kubernetes-resource/cmd/check
RUN go build -o /assets/in github.com/kdihalas/concourse-do-kubernetes-resource/cmd/in
RUN go build -o /assets/out github.com/kdihalas/concourse-do-kubernetes-resource/cmd/out

FROM gcr.io/distroless/static-debian11 AS resource
COPY --from=builder /assets /opt/resource