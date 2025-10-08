#!/bin/bash
set -e

V=$1

if [ -z "$V" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

git tag -a "$V" -m "$V"
git push origin "$V"

# GHCR
docker build -t ghcr.io/myferr/deo:"$V" .
docker tag ghcr.io/myferr/deo:"$V" ghcr.io/myferr/deo:latest
docker push ghcr.io/myferr/deo:"$V"
docker push ghcr.io/myferr/deo:latest

# Docker Hub
docker build -t myferr/deo:"$V" .
docker tag myferr/deo:"$V" myferr/deo:latest
docker push myferr/deo:"$V"
docker push myferr/deo:latest
