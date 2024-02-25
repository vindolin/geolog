FROM golang:1.22.0-bookworm

ARG ACCOUNT_ID
ARG LICENSE_KEY

RUN : "${ACCOUNT_ID:?You must provide a MaxMind account ID.}"
RUN : "${LICENSE_KEY:?You must provide a MaxMind license key.}"

WORKDIR $GOPATH/src/github.com/vindolin/geolog
RUN set -ex \
  && GEOIP_ARCHIVE_URL=https://$ACCOUNT_ID:$LICENSE_KEY@download.maxmind.com/geoip/databases/GeoLite2-City/download \
  && SHA256=$(wget -O - $GEOIP_ARCHIVE_URL?suffix=tar.gz.sha256 2>/dev/null | grep -o "^\w*\b") \
  && wget -O /tmp/GeoLite2-City.tar.gz $GEOIP_ARCHIVE_URL?suffix=tar.gz \
  && echo "$SHA256 /tmp/GeoLite2-City.tar.gz" | sha256sum -c -  \
  && MMDB_FILE=$(tar -ztf /tmp/GeoLite2-City.tar.gz | grep GeoLite2-City.mmdb) \
  && tar -zxf /tmp/GeoLite2-City.tar.gz $MMDB_FILE -O > /GeoLite2-City.mmdb

COPY geolog/*.go ./
COPY geolog/go.* ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /geolog

COPY geolog/index.html index.html
COPY geolog/favicon.ico favicon.ico

COPY run.sh run.sh
CMD ./run.sh
