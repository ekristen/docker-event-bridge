FROM scratch

ENTRYPOINT ["/bridge"]

COPY release/docker-event-bridge_linux-386 /bridge

