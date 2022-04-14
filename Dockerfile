ARG ENVOY_VERSION=v1.21-latest

FROM golang:1.16-alpine as build

WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY pkg /src/pkg
COPY Makefile Makefile
RUN apk update \
    && apk add make \
    && make build

FROM envoyproxy/envoy-alpine:${ENVOY_VERSION}

ENV CLUSTER_INFO=envoy-cluster
ENV type="\"@type\""

RUN mkdir -p /etc/envoy/tpl/xds/cds \
    mkdir -p /etc/envoy/tpl/xds/lds \
    mkdir -p /etc/envoy/tpl/xds/anchor \
    mkdir -p /etc/envoy/config

COPY --chown=envoy:envoy --from=build /src/bin/yaml-merge /bin/yaml-merge
COPY --chown=envoy:envoy docker-entorypoint.sh /bin/docker-entorypoint.sh
COPY --chown=envoy:envoy envoy.yaml /etc/envoy/tpl/envoy.yaml

ENTRYPOINT [ "/bin/docker-entorypoint.sh" ]

CMD [ "envoy", "-c", "/etc/envoy/envoy.yaml" ]
