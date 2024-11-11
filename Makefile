# e2e 测试
.PHONY: e2e_up
e2e_up:
	docker compose -f script/docker-compose.yml up -d

.PHONY: e2e_down
e2e_down:
	docker compose -f script/docker-compose.yml down

.PHONY: mock
mock:
	mockgen -package=mocks -destination=mocks/redis_cmdable.mock.go github.com/redis/go-redis/v9 Cmdable