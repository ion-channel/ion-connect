FROM golang:1.10.3-alpine3.8

RUN apk update && \
    apk upgrade && \
    apk add \
      git && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/ion-channel/ion-connect/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
RUN go get github.com/GeertJohan/go.rice
RUN go get github.com/GeertJohan/go.rice/rice
RUN rice append --exec ion-connect  -i ./lib

FROM scratch

ARG APP_NAME
ARG BUILD_DATE
ARG VERSION
ARG GIT_COMMIT_HASH
ARG ENVIRONMENT

LABEL org.metadata.build-date=$BUILD_DATE \
      org.metadata.version=$VERSION \
      org.metadata.vcs-url="https://github.com/ion-channel/ion-connect" \
      org.metadata.vcs-commit-id=$GIT_COMMIT_HASH \
      org.metadata.name="Ion Connect" \
      org.metadata.description="Ion Channel API Tool"

WORKDIR /root/

COPY --from=0 /etc/ssl /etc/ssl
COPY --from=0 /go/src/github.com/ion-channel/ion-connect/ion-connect .

ENTRYPOINT ["./ion-connect"]
