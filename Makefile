run:
	APP_NAME=mytheresa \
    MYTHERESA_MYSQL_HOST=127.0.0.1 \
    MYTHERESA_MYSQL_PORT=3306 \
    MYTHERESA_MYSQL_USER=mytheresa \
    MYTHERESA_MYSQL_PASSWORD=mytheresa \
    MYTHERESA_MYSQL_DB=mytheresa \
    MYTHERESA_REDIS_HOST=127.0.0.1 \
    MYTHERESA_REDIS_PORT=6379 \
    MYTHERESA_REDIS_PASSWORD= \
    HASH_ID_SAlT=mytheresa-salt-value \
    go run cmd/mytheresa/main.go
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
	docker build --progress=plain -t ginktest -f deploy/docker/mytheresa/Dockerfile-test .
	docker run --rm -it ginktest
