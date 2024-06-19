.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f lcsin/webook:v0.0.1
	@docker build -t lcsin/webook:v0.0.1 .

run:
	@go build -o webook.exe
	@./webook.exe

wire:
	@wire

mock:
	# service mock
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
	# repository mock
	@mockgen -source=internal/repository/user.go -package=repomocks -destination=internal/repository/mocks/user.mock.go
	@mockgen -source=internal/repository/code.go -package=repomocks -destination=internal/repository/mocks/code.mock.go
	# dao/cache mock
	@mockgen -source=internal/repository/dao/user.go -package=daomocks -destination=internal/repository/dao/mocks/user.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=cachemocks -destination=internal/repository/cache/mocks/user.mock.go
	# redis mock
	@mockgen -package=redismocks -destination=internal/repository/cache/redismock/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable

	@go mod tidy