run:
	@TEST_CONCURRENCY=2 TEST_HOST="www.baidu.com" TEST_ENV="dev" go run *.go

docker-build:
	@docker build -t yunzhu/httptest .

docker-push:
	@docker push yunzhu/httptest
