run:
	@TEST_CONCURRENCY=2 TEST_DEFAULT_ADDRESS="www.baidu.com" go run *.go

docker-build:
	@docker build -t yunzhu/httptest .

docker-push:
	@docker push yunzhu/httptest
