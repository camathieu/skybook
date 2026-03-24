.PHONY: all frontend server dev clean lint test test-frontend

# Build everything
all: frontend server

# Build Vue webapp → webapp/dist/
frontend:
	cd webapp && npm ci && npm run build

# Build Go server binary → server/skybook
# Copies webapp/dist into server/server/dist for go:embed, then builds.
server:
	rm -rf server/server/dist
	cp -r webapp/dist server/server/dist
	cd server && go build -o skybook .

# Dev mode: Vite on :5173 (proxies /api to :8080) + Go on :8080
dev:
	@echo "Starting Go server on :8080..."
	@cd server && go run . &
	@echo "Starting Vite dev server on :5173..."
	@cd webapp && npm run dev

# Lint Go code
lint:
	cd server && go fmt ./... && go vet ./...

# Go tests
test:
	cd server && go test ./... -count=1

# Frontend unit tests (vitest)
test-frontend:
	cd webapp && npm test

# Remove build artifacts
clean:
	rm -f server/skybook
	rm -rf server/server/dist
	rm -rf webapp/dist

# Docker
docker:
	docker build -t skybook .
