FROM index.docker.io/gesquive/go-builder:latest AS builder

ENV APP=cig
ARG TARGETARCH
ARG TARGETOS
ARG TARGETVARIANT

COPY dist/ /dist/
RUN copy-release

FROM scratch
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/${APP} /app/

# Use an unprivileged user
USER runner

ENTRYPOINT ["/app/cig"]
