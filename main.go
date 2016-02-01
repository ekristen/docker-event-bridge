package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
	consul "github.com/hashicorp/consul/api"
)

var (
	docker_client docker.Client
)

type ConsulEvent struct {
	Status string `json:"Status,omitempty" yaml:"Status,omitempty"`
	ID     string `json:"ID,omitempty" yaml:"ID,omitempty"`
	From   string `json:"From,omitempty" yaml:"From,omitempty"`
	Time   int64  `json:"Time,omitempty" yaml:"Time,omitempty"`
	Host   string `json:"Host,omitempty" yaml:"Host,omitempty"`
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func bridgeEvents(docker_socket string, consul_addr string, consul_token string, host string) {
	// Docker Client
	endpoint := "unix://" + docker_socket
	docker_client, err := docker.NewClient(endpoint)

	assert(err)

	consul_config := consul.DefaultConfig()
	consul_config.Address = consul_addr

	if consul_token != "" {
		consul_config.Token = consul_token
	}

	// Consul Client
	consul_client, err := consul.NewClient(consul_config)

	assert(err)

	consul_events := consul_client.Event()

	events := make(chan *docker.APIEvents)
	assert(docker_client.AddEventListener(events))
	log.Println("Listening for Docker events ...")

	quit := make(chan struct{})

	for msg := range events {
		b, json_err := json.Marshal(msg)
		if json_err != nil {
			log.Println("Docker Event: Error Parsing JSON", json_err)
			return
		}

		evt := ConsulEvent{}
		json.Unmarshal(b, &evt)
		evt.Host = host

		c, json_err := json.Marshal(evt)
		if json_err != nil {
			log.Println("Docker Event: Error Parsing JSON", json_err)
			return
		}

		e := &consul.UserEvent{Name: "docker", Payload: c}
		q := &consul.WriteOptions{}

		msg, _, fire_err := consul_events.Fire(e, q)
		if fire_err != nil {
			log.Println("Docker Event: Consul Event Fire Error,", fire_err)
			return
		}

		log.Println("Docker Event: Consul Event Delivered,", msg)
	}

	close(quit)
	log.Fatal("Docker event loop closed") // todo: reconnect?
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "docker-socket",
			Value: "/var/run/docker.sock",
			Usage: "Path to Docker Socket",
		},
		cli.StringFlag{
			Name:   "consul-http-addr",
			Value:  "localhost:8500",
			Usage:  "Consul HTTP ADDR",
			EnvVar: "CONSUL_ADDR",
		},
		cli.StringFlag{
			Name:   "consul-token",
			Usage:  "Consul ACL Token",
			EnvVar: "CONSUL_TOKEN",
		},
		cli.StringFlag{
			Name:   "hostname",
			Usage:  "Hostname to use for events",
			EnvVar: "HOSTNAME",
		},
	}

	app.Name = "consul-docker-event-bridge"
	app.Usage = "Send all docker events to consul"
	app.Version = "1.0.1"
	app.Action = func(c *cli.Context) {
		bridgeEvents(c.String("docker-socket"), c.String("consul-http-addr"), c.String("consul-token"), c.String("hostname"))
	}

	app.Run(os.Args)
}
