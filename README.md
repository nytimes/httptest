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
    -v $(pwd)/tests.yaml:/tests/tests.yaml \ # Mount test.yaml under /tests/ in container
    -e "TEST_HOST=example.com" \             # Specify hostname to test against
    yunzhu/httptest
```

You should see an output similar to this:
```
passed: root | [/]
1/1 tests passed
```

By default, the program parses all files in `$(pwd)/tests/` (in the above
example, `pwd` is `/`). This can be changed using an environment variable.

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

- GitHub Actions

  ```hcl
  action "httptest" {
    uses = "docker://yunzhu/httptest"
    env = {
      TEST_HOST = "example.com"
    }
  }
  ```

### Full test example

Any fields not explicitly stated as required are optional.

```yaml
tests:
  - description: 'root'          # Description, will be printed with test results (required)
    conditions:                  # Specify conditions. Test only runs when all conditions are met
      env:                       # Matches an environment variable
        TEST_ENV: '^(dev|stg)$'  # Environment variable name : regular expression
    request:                     # Request to send
      scheme: 'https'            # URL scheme. Only http and https are supported. Default: https
      host: 'example.com'        # Host to test against (this overrides TEST_HOST for this specific test)
      method: 'POST'             # HTTP method. Default: GET
      path: '/'                  # Path to hit. Default: /
      body: ''                   # Request body. Processed as string
    response:                    # Expected response
      statusCodes: [201]         # List of expected response status codes
      headers:                   # Expected response headers
        patterns:                # Match response header patterns
          server: '^BWS'         # Header name : regular expression
        notPresent:              # Specify headers not expected to exist.
          - 'abcdefg'
      body:                      # Response body
        patterns:                # Response body has to match all patterns in this list in order to pass test
          - '14px.*?"宋体"'

  - description: 'sign up page'  # Second test
    request:
      path: '/signup'
    response:
      statusCodes: [200]
```
