client:
	cd client && cargo build --release
server:
	cd server && go build ./...
test:
	cd server && go test ./... && go vet ./...
	cd client && cargo test && cargo clippy -- -D warnings
