FROM scratch

ENTRYPOINT ["/bridge"]

COPY dist/consul-docker-event-bridge_linux-386 /bridge

