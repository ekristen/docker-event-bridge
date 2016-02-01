FROM lalyos/scratch-chmx

ENTRYPOINT ["/bin/docker-event-bridge"]

ENV VERSION v1.0.0

ADD https://github.com/ekristen/docker-event-bridge/releases/download/$VERSION/docker-event-bridge_linux-386 /bin/docker-event-bridge

RUN ["/bin/chmx", "/bin/docker-event-bridge"]

