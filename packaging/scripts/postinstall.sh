#!/bin/sh
set -e

chown -R scw-exporter:scw-exporter /var/lib/scw-exporter
chmod 750 /var/lib/scw-exporter

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet scw-exporter.service; then
        systemctl restart scw-exporter.service
    fi
fi
