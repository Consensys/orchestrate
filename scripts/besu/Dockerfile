

FROM alpine:latest AS builder

RUN apk --no-cache add \
    curl \
    tar

RUN mkdir -p /pegasys-plus/plugins

ARG BINTRAY_USER_NAME
ARG BINTRAY_PASSWORD
ARG PPLUS_VERSION

# Download PegaSys Plus Plugins
RUN curl -u $BINTRAY_USER_NAME:$BINTRAY_PASSWORD -L https://dl.bintray.com/consensys/pegasys-enterprise-repo/pegasys-plus-$PPLUS_VERSION.tar.gz -o pegasys-plus.tar.gz || echo "Could not download P+ plugins"
RUN tar zxvf pegasys-plus.tar.gz || true
RUN mv /pegasys-plus-$PPLUS_VERSION/plugins/* /pegasys-plus/plugins || true

FROM hyperledger/besu:develop

COPY --from=builder /pegasys-plus/plugins /opt/besu/pegasys-plus/plugins

ENV BESU_OPTS="-Dbesu.plugins.dir=/opt/besu/pegasys-plus/plugins"