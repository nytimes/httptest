v# httptest
A simple concurrent HTTP API testing tool

## Usage

### Write a simple test

Create a new directory `tests` then create file `test.yaml` under this
directory with the following content:

```yaml
tests:
  - description: 'root'  # Description, will be printed with test results
    request:             # Request definition
      path: '/'          # Path
    response:            # Expected response
      statusCodes: [200] # List of expected response status codes
```

### Run tests locally

This program is distributed as a Docker image. By default, it parses all files
in `$(pwd)/tests/`. This can be changed using an environment variable.

```
docker run --rm \
    -v $(pwd)/tests:/tests \
    -e "TEST_HOST=example.com" \
    yunzhu/httptest
```

You should see an output similar to this
```
passed: root | [/]
1/1 tests passed
```

### Run tests in a CI/CD pipeline

This image can be used in any CI/CD system that supports Docker containers.

Examples

- Drone

  ```yaml
  pipeline:
    tests:
      image: yunzhu/httptest
      pull: true
      environment:
        TEST_HOST: 'example.com'
  ```

- Drone

  ```yaml
  pipeline:
    tests:
      image: yunzhu/httptest
      pull: true
      environment:
        TEST_HOST: 'example.com'
  ```


```yaml
tests:
  - description: 'baidu hp'
    conditions:
      env:
        TEST_ENV: '^(stg|prd)$'
    request:
      address: 'www.baidu.com'
      path: '/'
    response:
      statusCodes: [201]
      headers:
        patterns:
          server: '^BWS'
        notPresent:
          - 'abcdefg'
      body:
        patterns:
          - '14px.*?"宋体"'

  - description: 'hp'
    conditions:
      env:
        TEST_ENV: '^(stg|prd)$'
    request:
      path: '/'
    response:
      statusCodes: [201]

  - description: 'hp'
    request:
      path: '/'
    response:
      statusCodes: [200]
```
