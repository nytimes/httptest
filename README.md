# httptest
A simple concurrent HTTP testing tool

## Usage

### Write a simple test

Create a file `test.yaml` with the following content:

```yaml
tests:
  - description: 'root'  # Description, will be printed with test results
    request:             # Request to send
      path: '/'          # Path
    response:            # Expected response
      statusCodes: [200] # List of expected response status codes
```

### Run tests locally

This program is distributed as a Docker image. To run a container locally:
```bash
docker run --rm \
    -v $(pwd)/tests.yaml:/tests/tests.yaml \
    -e "TEST_HOST=example.com" \
    blupig/httptest
```

You should see an output similar to this:
```
passed:  tests.yaml | root | [/]

1 passed
0 failed
```

By default, the program parses all files in `$(pwd)/tests/` (in the above
example, `pwd` is `/` in the container). This can be changed using an
environment variable.

### Run tests in a CI/CD pipeline

This image can be used in any CI/CD system that supports Docker containers.

Examples

- Drone

  ```yaml
  pipeline:
    tests:
      image: blupig/httptest
      pull: true
      environment:
        TEST_HOST: 'example.com'
  ```

- GitHub Actions

  ```hcl
  action "httptest" {
    uses = "docker://blupig/httptest"
    env = {
      TEST_HOST = "example.com"
    }
  }
  ```

### Configurations

Environment variables:

- `TEST_HOST`: Host to test. Can be overridden by `request.host` in YAML.
  At least one of them needs to be specified, otherwise tests will fail.

- `TEST_CONCURRENCY`: Maximum number of requests can be sent at a time.
  Default: `2`.

- `TEST_DNS_OVERRIDE`: Override the IP address for `TEST_HOST`. Does not work
  for `request.host` specified in YAML.

- `TEST_PRINT_FAILED_ONLY`: Only print failed tests.
  Valid values: `false` or `true`. Default: `false`.


### Full test example

Required fields:

- `description`
- `request.path`

All other fields are optional. All matchings are case insensitive.

```yaml
tests:
  - description: 'root'          # Description, will be printed with test results. Required
    conditions:                  # Specify conditions. Test only runs when all conditions are met
      env:                       # Matches an environment variable
        TEST_ENV: '^(dev|stg)$'  # Environment variable name : regular expression

    request:                     # Request to send
      scheme: 'https'            # URL scheme. Only http and https are supported. Default: https
      host: 'example.com'        # Host to test against (this overrides TEST_HOST for this specific test)
      method: 'POST'             # HTTP method. Default: GET
      path: '/'                  # Path to hit. Required
      headers:                   # Headers
        x-test-header-0: 'abc'
        x-test-header-1: 'def'
      body: ''                   # Request body. Processed as string

    response:                    # Expected response
      statusCodes: [201]         # List of expected response status codes
      headers:                   # Expected response headers
        patterns:                # Match response header patterns
          server: '^ECS$'        # Header name : regular expression
          cache-control: '.+'
        notPresent:              # Specify headers not expected to exist.
          - 'set-cookie'         # Not regular expressions
          - 'x-frame-options'
      body:                      # Response body
        patterns:                # Response body has to match all patterns in this list in order to pass test
          - 'charset="utf-8"'    # Regular expressions
          - 'Example Domain'

  - description: 'sign up page'  # Second test
    request:
      path: '/signup'
    response:
      statusCodes: [200]
```
