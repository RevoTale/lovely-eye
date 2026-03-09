FROM mcr.microsoft.com/devcontainers/base:2-bookworm

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install additional tools
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y postgresql-client fish \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
