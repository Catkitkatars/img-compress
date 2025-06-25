FROM golang:1.24 AS build
WORKDIR /src
COPY . .
RUN --mount=type=ssh --mount=type=cache,target="/go/pkg/mod" \
    go mod download
RUN --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=cache,target="/go/pkg/mod" \
    CGO_ENABLED=1 go build -o img-compress cmd/img-compress/main.go

FROM ubuntu:22.04
WORKDIR /img-compress
RUN apt-get update && apt-get install -y \
    libsqlite3-0 && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /img-compress/storage

COPY --from=build /src/img-compress /img-compress/img-compress
COPY --from=build /src/config/local.yaml /img-compress/config/local.yaml

EXPOSE 8080
ENTRYPOINT ["./img-compress"]