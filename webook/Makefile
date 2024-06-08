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