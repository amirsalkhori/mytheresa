up:
	docker-compose up -d;
build:
	docker-compose up --build;
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

seed-products:
	docker-compose exec app /app/seeder 
	# @go run ./cmd/seeder/main.go  --products=$(PRODUCTS) --discounts=$(DISCOUNTS)

test-ginkgo:
	docker build --progress=plain -t ginktest -f deploy/docker/mytheresa/test.Dockerfile .
	docker run --rm -it ginktest
