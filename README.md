# httptest

[![Build Status](https://cloud.drone.io/api/badges/nytimes/httptest/status.svg)](https://cloud.drone.io/nytimes/httptest)

A simple concurrent HTTP testing tool

## Usage

### Write a simple test

Create a file `tests.yaml` with the following content:

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
    nytimes/httptest
```

You should see an output similar to this:
```
passed:  tests.yaml | root | /

1 passed
0 failed
0 skipped
```

Tip: If your test cases have conditions on environment variables (see `conditions` in [full example](#full-test-example)), remember to include `-e "<ENV_VAR>=<value>"`. e.g.

```
docker run --rm \
    -v $(pwd)/tests.yaml:/tests/tests.yaml \
    -e "TEST_HOST=stg.example.com" \
    -e "TEST_ENV=stg" \
    nytimes/httptest
```

By default, the program parses all files in `$(pwd)/tests` recursively.
This can be changed using an environment variable.

### Run tests in a CI/CD pipeline

This image can be used in any CI/CD system that supports Docker containers.

Examples

- Drone

  ```yaml
  pipeline:
    tests:
      image: nytimes/httptest
      pull: true
      environment:
        TEST_HOST: 'example.com'
  ```

- GitHub Actions

  ```hcl
  action "httptest" {
    uses = "docker://nytimes/httptest"
    env = {
      TEST_HOST = "example.com"
    }
  }
  ```

### Configurations

A few global configurations (applied to all tests) can be specified by
environment variables:

- `TEST_DIRECTORY`: Local directory that contains the test definition YAML
  files. Default: `tests`

- `TEST_HOST`: Host to test. Can be overridden by `request.host` of individual
  test definitions. If `TEST_HOST` and `request.host` are both not set, test
  will fail.

- `TEST_CONCURRENCY`: Maximum number of concurrent requests at a time.
  Default: `2`.

- `TEST_DNS_OVERRIDE`: Override the IP address for `TEST_HOST`. Does not work
  for `request.host` specified in YAML.

- `TEST_PRINT_FAILED_ONLY`: Only print failed tests. Valid values: `false` or
  `true`. Default: `false`.

### Environment variable substitution

This program supports variable substitution from environment variables in YAML
files. This is useful for handling secrets or dynamic values.
Visit [here](https://github.com/drone/envsubst/blob/master/README) for
supported functions.

Example:

```yaml
tests:
  - description: 'get current user'
    request:
      path: '/user'
      headers:
        authorization: 'token ${SECRET_AUTH_TOKEN}'
    response:
      statusCodes: [200]
```

### Full test example

Required fields for each test:

- `description`
- `request.path`

All other fields are optional. All matchings are case insensitive.

```yaml
tests:
  - description: 'root'          # Description, will be printed with test results. Required
    conditions:                  # Specify conditions. Test only runs when all conditions are met
      env:                       # Matches an environment variable
        TEST_ENV: '^(dev|stg)$'  # Environment variable name : regular expression
    skipCertVerification: false  # Set true to skip verification of server TLS certificate (insecure and not recommended)

    request:                     # Request to send
      scheme: 'https'            # URL scheme. Only http and https are supported. Default: https
      host: 'example.com'        # Host to test against (this overrides TEST_HOST for this specific test)
      method: 'POST'             # HTTP method. Default: GET
      path: '/'                  # Path to hit. Required
      headers:                   # Headers
        x-test-header-0: 'abc'
        x-test: '${REQ_TEST}'    # Environment variable substitution
      body: ''                   # Request body. Processed as string

    response:                    # Expected response
      statusCodes: [201]         # List of expected response status codes
      headers:                   # Expected response headers
        patterns:                # Match response header patterns
          server: '^ECS$'        # Header name : regular expression
          cache-control: '.+'
        notPresent:              # Specify headers not expected to exist.
          - 'set-cookie'         # These are not regular expressions
          - 'x-frame-options'
        notMatching:
          set-cookie: ^.*abc.*$  # Specify headers expected to exist but NOT match the given regex
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

## Development

### Run locally

Download package
```
go get -d -u github.com/nytimes/httptest
```

Build and run
```bash
# In repo root directory
make run
```
This will run the tests defined in `example-tests` directory.
