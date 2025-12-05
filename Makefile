.PHONY: client
client:
	(cd client && go run main.go);

.PHONY: server
server:
	(cd server && go run main.go);

.PHONY: test
test:
	(cd client && go test ./...);
	(cd server && go test ./...);

.PHONY: fmt
fmt:
	(cd client && go fmt ./...);
	(cd server && go fmt ./...);

.PHONY: e2e
e2e:
	@echo "Starting server..."
	@(cd server && go run main.go) &
	@sleep 2
	@echo "Running e2e tests..."
	@(cd client && go test -run TestE2ERoomJoining -v)
	@echo "Stopping server..."
	@pkill -f "go run main.go" || true

.PHONY: clean
clean:
	rm -rf client/client server/server
