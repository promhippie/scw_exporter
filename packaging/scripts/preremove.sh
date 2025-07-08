#!/bin/sh
set -e

systemctl stop scw-exporter.service || true
systemctl disable scw-exporter.service || true
