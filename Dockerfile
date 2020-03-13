FROM scratch

ARG GIT_COMMIT_HASH

LABEL org.metadata.vcs-url="https://github.com/ion-channel/ion-connect" \
      org.metadata.vcs-commit-id=$GIT_COMMIT_HASH \
      org.metadata.name="Ion Connect" \
      org.metadata.description="Ion Channel CLI Tool"

COPY --from=alpine /etc/ssl /etc/ssl
ADD ion-connect /ion-connect

WORKDIR /data

CMD ["/ion-connect", "-v"]
