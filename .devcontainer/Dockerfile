FROM mcr.microsoft.com/devcontainers/go:2.2-1.26-trixie@sha256:9dd796c5a81000d4e507d7002b36d3eae0e6f5bbc9eaa49afe3e5a424534ed15

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Install additional tools
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y postgresql-client fish \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
