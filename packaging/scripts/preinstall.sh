#!/bin/sh
set -e

if ! getent group scw-exporter >/dev/null 2>&1; then
    groupadd --system scw-exporter
fi

if ! getent passwd scw-exporter >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/scw-exporter --shell /bin/bash -g scw-exporter scw-exporter
fi
