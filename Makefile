
test:
	ginkgo -r
up:
	docker-compose up -d;
build-up:
	docker-compose up --build -d;
ps:
	docker-compose ps;
down:
	docker-compose down;
build-up:
	docker-compose up -d --force-recreate --build;
verbose-build-up:
	docker-compose up --force-recreate --build;
container-logs:
	docker-compose logs -f;
env:
	@[ -e ./.env ] || cp -v ./.env.example ./.env
tidy:
	go mod tidy

test-repository:
	TEST_TYPE=integration go test ./internal/infra/db/mysql/repository... -v
seed-products:
	@go run ./cmd/seeder/main.go  --products=$(PRODUCTS) --discounts=$(DISCOUNTS)
