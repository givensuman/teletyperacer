.PHONY: build server client multiplayer

# Will need a better way to handle this 
# in the future
TERM = cosmic-term

build:
	docker compose build

.PHONY: server
server: # Spin up the server
	docker compose up server

.PHONY: client
client: # Spin up the client
	docker compose run --rm client

.PHONY: multiplayer
multiplayer: # Spin up two clients and one server
	docker compose up -d server
	sleep 2
	$(TERM) -e docker compose run --rm client &
	$(TERM) -e docker compose run --rm client &
	wait
