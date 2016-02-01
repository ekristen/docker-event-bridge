FROM scratch

ENTRYPOINT ["/bridge"]

COPY dist/docker-event-bridge_linux-386 /bridge

