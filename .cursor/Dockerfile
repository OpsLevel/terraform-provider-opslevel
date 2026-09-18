# Environment image for Cursor cloud/background agents working on
# terraform-provider-opslevel.
#
# Referenced by .cursor/environment.json. Cursor builds this image, then clones
# the repo and runs the agent inside it as the "ubuntu" user, so this file only
# installs the toolchain -- it never copies the project. Dependency warming
# happens in the "install" script in environment.json, which runs once per Build
# with the workspace present.

FROM golang:1.24-bookworm

# Versions are pinned to match .github/workflows; override with --build-arg.
ARG TERRAFORM_VERSION=1.7.0
ARG TASK_VERSION=3.40.0
ARG GOLANGCI_LINT_VERSION=1.64.8
# Newest releases that still build with Go 1.24 (GOTOOLCHAIN=local). Current
# @latest for both requires Go >= 1.26, and golangci-lint @latest is v2, which
# rejects this repo's v1-format .golangci.yml.
ARG GOFUMPT_VERSION=0.9.1
ARG CHANGIE_VERSION=1.22.0

ENV GOTOOLCHAIN=local \
    CHECKPOINT_DISABLE=1 \
    TF_IN_AUTOMATION=1

# sudo and git are required by Cursor agents.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates curl git jq unzip sudo openssh-client \
    && rm -rf /var/lib/apt/lists/*

# Cursor runs the agent as a non-root "ubuntu" user at UID/GID 1000. The Go
# images are Debian-based and have no such user, so create it if absent.
RUN set -eux; \
    if ! getent group 1000 >/dev/null; then groupadd -g 1000 ubuntu; fi; \
    if ! id -u 1000 >/dev/null 2>&1; then \
        useradd -m -u 1000 -g 1000 -s /bin/bash ubuntu; \
    fi; \
    echo 'ubuntu ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/ubuntu; \
    chmod 0440 /etc/sudoers.d/ubuntu

# Terraform
RUN set -eux; \
    arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/terraform.zip \
        "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${arch}.zip"; \
    unzip -d /usr/local/bin /tmp/terraform.zip; \
    rm /tmp/terraform.zip; \
    terraform version

# Task runner
RUN set -eux; \
    arch="$(dpkg --print-architecture)"; \
    curl -fsSL -o /tmp/task.tar.gz \
        "https://github.com/go-task/task/releases/download/v${TASK_VERSION}/task_linux_${arch}.tar.gz"; \
    tar -C /usr/local/bin -xzf /tmp/task.tar.gz task; \
    rm /tmp/task.tar.gz; \
    task --version

# Go tooling used by `task lint` / `task ci`. Everything lands in /usr/local/bin
# rather than a Go-specific bin dir: Cursor resets PATH at session start, so only
# directories on the default PATH can be relied on.
RUN set -eux; \
    curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sh -s -- -b /usr/local/bin "v${GOLANGCI_LINT_VERSION}"; \
    GOBIN=/usr/local/bin go install "mvdan.cc/gofumpt@v${GOFUMPT_VERSION}"; \
    GOBIN=/usr/local/bin go install "github.com/miniscruff/changie@v${CHANGIE_VERSION}"; \
    golangci-lint --version

# Same reason: the base image puts go on PATH via ENV, which a PATH reset drops.
RUN ln -sf /usr/local/go/bin/go /usr/local/bin/go \
    && ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

USER ubuntu

# The base image sets GOPATH=/go, which is root-owned and unwritable by the
# agent user -- and an inherited ENV beats ~/.config/go/env, so it has to be
# overridden here as well as below.
ENV GOPATH=/home/ubuntu/go

# Written to ~/.config/go/env too, so the values still hold if Cursor resets the
# environment at session start. Caches live under $HOME so the agent can write them.
RUN go env -w \
        GOPATH=/home/ubuntu/go \
        GOMODCACHE=/home/ubuntu/go/pkg/mod \
        GOCACHE=/home/ubuntu/.cache/go-build \
        GOBIN=/home/ubuntu/go/bin \
        GOTOOLCHAIN=local

# Belt and braces for the PATH reset, per Cursor's guidance to persist tool
# visibility in the runtime user's shell profile.
RUN printf '\n# terraform-provider-opslevel toolchain\nexport PATH="/usr/local/bin:/usr/local/go/bin:/home/ubuntu/go/bin:$PATH"\n' \
        >> /home/ubuntu/.bashrc

# `task build` derives the provider version from `git describe`, which refuses a
# workspace owned by a different uid.
RUN git config --global --add safe.directory '*'

# .gitmodules uses git@github.com: URLs, which need SSH credentials the agent
# does not have. Both submodule repos are public, so rewrite to anonymous HTTPS.
# The submodules are initialized by the install script in environment.json.
RUN git config --global url."https://github.com/".insteadOf "git@github.com:"

# No CMD/ENTRYPOINT on purpose: Cursor controls the container lifecycle and runs
# the agent and terminals itself.
