.PHONY: debug
debug: docker-check
	@docker compose --profile dev down && docker system prune --volumes --force && docker compose --profile dev up -d

.PHONY: prod
prod: docker-check docker-build
	@docker compose --profile prod down && docker system prune --volumes --force && docker compose --profile prod up

.PHONY: test
test:
	@go test -count=1 -covermode=atomic ./... | grep -v "/generated/"

.PHONY: test-race
test-race:
	@go test -race -count=1

.PHONY: lint
lint:
	@go tool golangci-lint run

.PHONY: format
format:
	@go tool gofumpt -l -w . && go tool golines -w . && go tool goimports -w -local "avito_pvz/" .

.PHONY: db-cli
db-cli:
	@PGPASSWORD=password pgcli --host 127.0.0.1 --port 5432 --username postgres

# brew install grpc
# brew install protoc-gen-go
# brew install protoc-gen-go-grpc
.PHONY: codegen
codegen:
	@go tool oapi-codegen --config=.oapi-codegen.yaml assignment/swagger.yaml \
	&& protoc --proto_path=assignment assignment/*.proto --go_out=internal/generated/grpc --go-grpc_out=internal/generated/grpc \
	&& mockery --log-level="" && rm -rf internal/generated/mocks && mkdir internal/generated/mocks && mockery --log-level=""

.PHONY: docker-build
docker-build: docker-check
	@docker build --no-cache --rm --platform linux/amd64,linux/arm64 --tag avito-pvz-server .

.PHONY: docker-check
docker-check:
	@if ! command -v docker &> /dev/null; then \
		echo "Error: Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi

