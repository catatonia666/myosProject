
ARG GO_VERSION=1.24.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x


COPY . .

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/apiserver

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
WORKDIR /app    
COPY --from=build /bin/server /bin/
COPY --from=build /app/ui /app/ui
COPY --from=build /app/config /app/config
COPY --from=build /app/ui/static /app/ui/static
EXPOSE 3000

ENTRYPOINT [ "/bin/server" ]