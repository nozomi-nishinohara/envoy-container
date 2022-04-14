#!/usr/bin/env sh
set -e

# cds
yaml-merge -i /etc/envoy/tpl/xds/cds -d /etc/envoy/tpl/xds/anchor -o /etc/envoy/config/cds.yaml

# lds
yaml-merge -i /etc/envoy/tpl/xds/lds -d /etc/envoy/tpl/xds/anchor -o /etc/envoy/config/lds.yaml

# envoy.yaml
yaml-merge -i /etc/envoy/tpl/envoy.yaml -o /etc/envoy/envoy.yaml


. /docker-entrypoint.sh "$@"