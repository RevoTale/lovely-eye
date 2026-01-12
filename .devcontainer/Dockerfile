FROM mcr.microsoft.com/devcontainers/go:2-1.25-bookworm

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install additional tools
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y postgresql-client \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
