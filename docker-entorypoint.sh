#!/usr/bin/env sh
set -e

ANCHOR=${ANCHOR:-anchor}

for xds in `echo $XDS | tr -s ',' ' '`; do
    # cds
    yaml-merge -i "/etc/envoy/tpl/xds/$xds" -d "/etc/envoy/tpl/xds/$ANCHOR" -o "/etc/envoy/config/$xds.yaml"
done

# envoy.yaml
yaml-merge -i /etc/envoy/tpl/envoy.yaml -o /etc/envoy/envoy.yaml


. /docker-entrypoint.sh "$@"