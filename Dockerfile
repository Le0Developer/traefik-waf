FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

COPY . .

RUN GOPATH= GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags no_fs_access -o /main .

FROM scratch

COPY --from=builder main main

ENTRYPOINT ["/main"]
