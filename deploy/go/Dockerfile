FROM        scratch

ADD         ./app /app
WORKDIR     /data

ENV         GRPC_PORT 5353
ENV         WS_PORT 5354
EXPOSE      5353
EXPOSE      5354

ENTRYPOINT  ["/app"]
