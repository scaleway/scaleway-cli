# Scripts for Scaleway CLI development

| Script         | Description                                                              |
|----------------|--------------------------------------------------------------------------|
| `build.sh`     | Build the CLI binary for all supported platform (Linux, Darwin, Windows) |
| `install.sh`   | Build and install the CLI binary (MacOSX support only)                   |
| `lint.sh`      | Run the linter golangci-lint on the codebase                             |
| `run-tests.sh` | Execute the test suite for the CLI v2                                    |

## Requirements

Some of those scripts have external requirements such as:

- [go](https://golang.org/doc/install)
- [golangci-lint](https://github.com/golangci/golangci-lint#install)

## build.sh

```
$ ./scripts/build.sh   
```

## install.sh

```
./scripts/install.sh
```

## lint.sh

```
$ ./scripts/lint.sh -h
Usage:
  ./scripts/lint.sh [OPTIONS]

Options:
  -w, --write
	Fix found issues (if it's supported by the linter).
  --list
	List current linters configuration.
  -h, --help
	Display this help.

```

## run-tests.sh

CLI testing uses cassettes and golden files to run a part of its test suite.

### Cassettes

Cassettes record a set of interactions between the CLI and the Scaleway API.
A given test would define a scenario that would be described as a YAML file that contains all those interactions. 
For instance:

```
version: 1
interactions:
- request:
    body: ""
    form: {}
    headers:
      User-Agent:
      - scaleway-sdk-go/v1.0.0-beta.5+dev (go1.13.6; darwin; amd64) cli-e2e-test
    url: https://api.scaleway.com/
    method: GET
  response:
    body: '{"Name":"Scaleway Api","Description":"Welcome to the Scaleway public API!","Version":"v0.0.140","ProtobufVersion":"a0be3c28","DocumentationUrl":"https://developers.scaleway.com"}
```

### Golden

A golden file stores the output of a command and is used as a reference by the test as the expected output.
Those files should be updated any time the output changes for reasons such as: API changes, new format, etc.
The output is stored as a separate file rather than as a string literal inside the test code. 
So when the test is executed, it will read in the file and compare it to the output produced by the system under test.

### Help command

```
$ ./scripts/run-tests.sh -h
Usage:
  ./scripts/run-tests.sh [OPTIONS]

Options:
  -r, --run <regex>
	Run a specific test or set of tests matching the given regex. Similar to the '-run' Go test flag.
  -u, --update
	Update goldens and record cassettes during integration tests.
  -g, --update-goldens
	Update goldens during integration tests.
  -c, --update-cassettes
	Record cassettes during integration tests. Warning: a valid Scaleway token is required in your environment in order to record cassettes.
  -D, --debug
	Enable CLI debug mode.
  -h, --help
	Display this help.
```
