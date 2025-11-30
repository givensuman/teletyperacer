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

.PHONY: clean
clean:
	rm -rf client/client server/server
