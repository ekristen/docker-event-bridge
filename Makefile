build:
	@sh -c "'$(CURDIR)/scripts/build.sh'"
	docker build -t ekristen/consul-docker-event-bridge -f Dockerfile .

