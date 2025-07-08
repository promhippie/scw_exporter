#!/bin/sh
set -e

if [ ! -d /var/lib/scw-exporter ]; then
    userdel scw-exporter 2>/dev/null || true
    groupdel scw-exporter 2>/dev/null || true
fi
