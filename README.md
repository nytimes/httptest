# httptest

[![Build Status](https://cloud.drone.io/api/badges/nytimes/httptest/status.svg)](https://cloud.drone.io/nytimes/httptest)

A simple concurrent HTTP testing tool

## Usage

### Write a simple test

Create a file `tests.yml` with the following content:

```yml
tests:
  - description: 'root'  # Description, will be printed with test results
    request:             # Request to send
      path: '/'          # Path
    response:            # Expected response
      statusCodes: [200] # List of expected response status codes
```

### Run tests locally

This program is distributed as a Docker image. To run a container locally (under a folder contains tests related yml file): 
```bash
docker run --rm \
    -v $(pwd)/httpbin.yml:/tests/tests.yml \
    -e "TEST_HOST=example.com" \
    nytimes/httptest
```

You should see an output similar to this:
```
passed:  tests.yml | root | /

1 passed
0 failed
0 skipped
```

Tip: If your test cases have conditions on environment variables (see `conditions` in [full example](#full-test-example)), remember to include `-e "<ENV_VAR>=<value>"`. e.g.

```
docker run --rm \
    -v $(pwd)/route-abc-tests.yml:/tests/tests.yml \
    -e "TEST_HOST=www.stg.example.com" \
    -e "TEST_ENV=stg" \
    nytimes/httptest
```

By default, if your current working directory is `$(pwd)/tests`,  the program will parse all files in `$(pwd)/tests` recursively.
This can be changed using an environment variable.

### Run tests in a CI/CD pipeline

This image can be used in any CI/CD system that supports Docker containers.

Examples

- Drone

  ```yml
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

- `TEST_DIRECTORY`: Local directory that contains the test definition YML or YAML
  files. Default: `tests`

- `TEST_HOST`: Host to test. Can be overridden by `request.host` of individual
  test definitions. If `TEST_HOST` and `request.host` are both not set, test
  will fail.

- `TEST_CONCURRENCY`: Maximum number of concurrent requests at a time.
  Default: `2`.

- `TEST_DNS_OVERRIDE`: Override the IP address for `TEST_HOST`. Does not work
  for `request.host` specified in YML.

- `TEST_PRINT_FAILED_ONLY`: Only print failed tests. Valid values: `false` or
  `true`. Default: `false`.

- `TEST_VERBOSITY`: Increase logging output for tests. Default: `0`.

- `ENABLE_RETRIES`: Enables retrying requests if a test does not succeed.
   Defaults: `true`.

- `DEFAULT_RETRY_COUNT`: Specify the number of times to retry a test request if
   the initial request does not succeed. Only applied if `ENABLE_RETRIES` is set
   to `true` Defaults: `2`.

### Environment variable substitution

This program supports variable substitution from environment variables in YML
files. This is useful for handling secrets or dynamic values.
Visit [here](https://github.com/drone/envsubst) for
supported functions.

Example:

```yml
tests:
  - description: 'get current user'
    request:
      path: '/user'
      headers:
        authorization: 'token ${SECRET_AUTH_TOKEN}'
    response:
      statusCodes: [200]
```

### Dynamic headers

In addition to defining header values statically (i.e. when a test is written), there is also support for determining the value of a header at runtime. This feature is useful if the value of a header must be determined when a test runs (for example, a header that must contain a recent timestamp). Such headers are specified as `dynamicHeaders` in the YAML, and each must be accompanied by the name of a function to call and a list of arguments for that function (if required). For example:

```yaml
dynamicHeaders:
  - name: x-test-signature
    function: myFunction
    args:
      - 'arg1'
      - 'arg2'
      - 'arg3'
```

The arguments can be literal strings, environment variables, or values of previously set headers (see examples below). The function definitions are in the file `dynamic.go` and can be added to as necessary. Each function must implement the following function interface:

```go
type resolveHeader func(existingHeaders map[string]string, args []string) (string, error)
```

These are the functions that are currently supported:

#### now
Returns the number of seconds since the Unix epoch

Args: none

#### signStringRS256PKCS8
Constructs a string from args (delimited by newlines), signs it with the (possibly passphrase-encrypted) PKCS #8 private key, and returns the signature in base64

Args:
- key
- passphrase (can be empty string if key is not encrypted)
- string(s)... (any amount of strings, from previously set headers or literal values)

#### postFormURLEncoded
Sends an HTTP POST request to the given URL (with Content-Type `application/x-www-form-urlencoded`), returns either the whole response body or a specific JSON element, and treats the remaining arguments as data

Args:
- url (the URL to send the request to)
- element (see section on "JSON Query Syntax" below)
- string(s)... (any amount of strings, from previously set headers or literal values, which will be concatenated and delimited with '&' for the request body)

##### JSON Query Syntax
The element to return from the JSON response body is specified by a dot-delimited string representing the path to the element. Each part of this string is a single step in that path, and array access is handled by providing a zero-based index. If the element selected has a primitive value, that value is returned. If the element selected is a JSON object or array, that object/array is returned in compact form (insignificant whitespace removed). To return the entire JSON response body, use an empty string.

**Examples:**  
**Primitive data**  
Response body:
```json
{
    "nested": {
        "object": {
            "data": 123
        }
    }
}
```
Query: 'nested.object.data'  
Result: 123

**Object**  
Response body:
```json
{
    "nested": {
        "object": {
            "data": 123
        }
    }
}
```
Query: 'nested.object'  
Result: {"data":123}

**Array**  
Response body:
```json
{
    "nested": {
        "array": [
            "abc",
            "def",
            "ghi"
        ]
    }
}
```
Query: 'nested.array'  
Result: ["abc","def","ghi"]

**Array element**  
Response body:
```json
{
    "nested": {
        "array": [
            "abc",
            "def",
            "ghi"
        ]
    }
}
```
Query: 'nested.array.1'  
Result: def

**Whole response**  
Response body:
```
Just a normal, non-JSON response
```
Query: ''  
Result: Just a normal, non-JSON response

#### concat
Concatenates the arguments into a single string

Args:
- string(s)... (any amount of strings, from previously set headers or literal values, which will be concatenated)

### Full test example

Required fields for each test:

- `description`
- `request.path`

All other fields are optional. All matchings are case insensitive.

```yml
tests:
  - description: 'root'                        # Description, will be printed with test results. Required
    conditions:                                # Specify conditions. Test only runs when all conditions are met
      env:                                     # Matches an environment variable
        TEST_ENV: '^(dev|stg)$'                # Environment variable name : regular expression
    skipCertVerification: false                # Set true to skip verification of server TLS certificate (insecure and not recommended)

    request:                                   # Request to send
      scheme: 'https'                          # URL scheme. Only http and https are supported. Default: https
      host: 'example.com'                      # Host to test against (this overrides TEST_HOST for this specific test)
      method: 'POST'                           # HTTP method. Default: GET
      path: '/'                                # Path to hit. Required
      headers:                                 # Headers
        x-test-header-0: 'abc'
        x-test: '${REQ_TEST}'                  # Environment variable substitution
      dynamicHeaders:                          # Headers whose values are determined at runtime (see "Dynamic Headers" section above)
        - name: x-test-timestamp               # Name of the header to set
          function: now                        # Calling a no-arg function to get the header value
        - name: x-test-signature
          function: signStringRS256PKCS8       # Function called to determine the header value
          args:                                # List of arguments for the function (can be omitted for functions that don't take arguments)
            - '${TEST_KEY}'                    # Can use environment variable substitution
            - '${TEST_PASS}'
            - x-test-timestamp                 # Can use values of previously set headers
            - '/svc/login'                     # Can use literals
            - x-test-header-0
            - x-test
        - name: x-test-token
          function: postFormURLEncoded         # Function called to retrieve the header value from an API
          args:
            - '${TOKEN_URL}'                   # The URL to POST to
            - 'results.7.token.id_token'       # The dot-delimited string representing an element to return from the JSON response body; use integers for array elements or '' for the whole response body
            - 'client_secret=${CLIENT_SECRET}' # Can use literals, environment variables, or a combination
            - x-test-timestamp                 # Can use values of previously set headers
            - 'client_id=123'                  # These arguments will be combined to send 'client_secret=${CLIENT_SECRET}&x-test-timestamp&client_id=123' in the body of the request
        - name: authorization
          function: concat                     # Function that concatenates strings, either literals or values from previously set headers
          args:
            - 'Bearer '
            - x-test-token
      body: ''                                 # Request body. Processed as string

    response:                                  # Expected response
      statusCodes: [201]                       # List of expected response status codes
      headers:                                 # Expected response headers
        patterns:                              # Match response header patterns
          server: '^ECS$'                      # Header name : regular expression
          cache-control: '.+'
        notPresent:                            # Specify headers not expected to exist.
          - 'set-cookie'                       # These are not regular expressions
          - 'x-frame-options'
        notMatching:
          set-cookie: ^.*abc.*$                # Specify headers expected to exist but NOT match the given regex
      body:                                    # Response body
        patterns:                              # Response body has to match all patterns in this list in order to pass test
          - 'charset="utf-8"'                  # Regular expressions
          - 'Example Domain'

  - description: 'sign up page'                # Second test
    request:
      path: '/signup'
    response:
      statusCodes: [200]

  - description: 'HTTPS GET - test if present not matching'
    request:
      path: '/get'
    conditions:
      env:
        TEST_ENV: dev
    response:
      statusCodes: [200]
      headers:
        ifPresentNotMatching:
          Content-Type: ^notreal/fake$         # If Content-Type does not exist or does not match this pattern the test passes, else it fails

  - description: 'HTTPS GET - test if present not matching multiple possible matches'
    request:
      path: '/get'
    conditions:
      env:
        TEST_ENV: dev
    response:
      headers:
        notMatching: 
          Server: ^(.*(nginx|42)).*            # If Server does not match either portion of the pattern (nginx or 42) the test passes, else it fails
      statusCodes:
        - 200
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
