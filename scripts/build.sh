#!/bin/bash
set -e

if [ -z "$1" ]; then
    OS_PLATFORM_ARG=(-os="linux")
else
    OS_PLATFORM_ARG=($1)
fi

if [ -z "$2" ]; then
    OS_ARCH_ARG=(-arch="386")
else
    OS_ARCH_ARG=($2)
fi

# Build Docker image unless we opt out of it
if [[ -z "$SKIP_BUILD" ]]; then
   docker build --rm=true --force-rm=true -t deb-builder -f Dockerfile.build .
fi

# Get rid of existing binaries
rm -f *-386
rm -f *-amd64
rm -f release/*
docker run --rm -v `pwd`:/go/src/github.com/ekristen/docker-event-bridge:Z deb-builder gox "${OS_PLATFORM_ARG[@]}" "${OS_ARCH_ARG[@]}" -output="release/{{.Dir}}_{{.OS}}-{{.Arch}}" -ldflags="-w"
