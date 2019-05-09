run:
	@TEST_DIRECTORY="example-tests" TEST_CONCURRENCY=2 TEST_HOST="httpbin.org" TEST_ENV="dev" go run *.go
