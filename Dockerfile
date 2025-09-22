# -----------------------
# Build Stage
# -----------------------
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ubuntu-autoinstaller main.go

# -----------------------
# Runtime Stage (Ubuntu)
# -----------------------

# Version args used for image naming and labels
ARG UBUNTU_VERSION
ARG APP_VERSION=1.0.0
ARG RELEASE_TAG

# Use a default in the FROM line via shell-like fallback (requires BuildKit)
FROM ubuntu:${UBUNTU_VERSION:-22.04}


LABEL org.opencontainers.image.version=${APP_VERSION} \
      org.opencontainers.image.revision=${RELEASE_TAG} \
      org.opencontainers.image.base.name="ubuntu:${UBUNTU_VERSION}"

ENV UAI_VERSION=${APP_VERSION}
ENV UAI_PORT=8080
ENV UAI_MODE=release
ENV STATIC_DIR=/app/static

# Install required packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    curl \
    tzdata \
    tini \
    iputils-ping \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/ubuntu-autoinstaller .
COPY --from=builder /app/static ${STATIC_DIR}

EXPOSE ${UAI_PORT}

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://0.0.0.0:${UAI_PORT}/health || exit 1

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["./ubuntu-autoinstaller"]