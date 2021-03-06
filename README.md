# Docker Event Bridge

With this first version, it simply pushes all docker events up to a consul server. It would be nice to modularize the backends to support other destinations, but this helps roll up all events to central location if you are wanting to do something with the events globally across your infrastructure.

## Configuration

* `--consul-http-addr` - Specify where to connect to consul
* `--consul-token` - If you are using consul ACLs
* `--host` - This hostname will be attached to all events
* `--docker-socket` - Default: /var/run/docker.sock - Specifies where to talk to docker, currently only the socket is supported.

## Installation

`docker pull ekristen/docker-event-bridge`

## Running

```
docker run -d \
  -h $HOSTNAME \
  --name="docker-event-bridge" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ekristen/docker-event-bridge --consul-http-addr=10.10.10.10:8500 --consul-token=ABC123`
```

## Building

`make scripts/build.sh`

By default only linux builds are made.

You will need a local docker instance running that supports mounting in your host volume path into the container. If you are on OSX, this can be achieved using docker-machine (with virtualbox or vmware fusion). If you are on linux, and running a local docker daemon this is already supported.

## Credits

Borrowed ideas from https://github.com/gliderlabs/registrator, thanks GliderLabs!
