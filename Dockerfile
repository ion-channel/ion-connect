FROM alpine:3.5

MAINTAINER dev@ionchannel.io

ARG BUILD_DATE
ARG VERSION

LABEL org.metadata.base.build-date=$BUILD_DATE \
      org.metadata.base.version=$VERSION \
      org.metadata.name="Ion Channel Alpine ion-connect Image" \
      org.metadata.description="A base docker image for Ion Channel's ion-connect utility" \
      org.metadata.url="https://ionchannel.io" \
      org.metadata.vcs-url="https://github.com/ion-channel/ion-connect`"

RUN apk update && \
    apk upgrade && \
    apk add \
      bash jq wget && \
    rm -rf /var/cache/apk/*

COPY ion-connect /usr/bin/ion-connect

CMD ["ion-connect", "-v"]
