# Testing the CLI

TL;DR

| go test flags | Description                              |
|---------------|------------------------------------------|
| `-debug`      | Enable debug mode for the test           |
| `-goldens`    | Update the golden files of the run tests |
| `-cassettes`  | Update the cassettes of the run tests    |

### Objectives of the test suite

- Ensure that we got no new regression when we merge new code
- Avoid leaking credentials when doing integration testing with Scaleway APIs

### Solution

#### Interaction Recording: Golden & Cassettes files

A **cassette** file contains all the interactions of the Scaleway APIs that a CLI must produce when a given test run.
A **golden** file contains the output that the CLI must produce when a given interaction recorded by a cassette occurs.

Having golden ensure that any change in a command output will be noticed if the behavior of the CLI change for a given interaction.
Having cassette ensure that we can replay interactions without actually doing the calls.

Warning, if you choose to record a new cassette, you will create real resources on your organization and will be billed accordingly.
Be sure to check out the resource you create in each test and remember to delete them once you need them anymore.

#### Metadata

When running a test, you might need information such as ID that you cannot know in advance (such as ID of resources).
The `core.Test` uses different helpers to pass useful information around.

One of them is the [`core.testMetadata`](https://github.com/scaleway/scaleway-cli/blob/master/internal/core/testing.go#L80).
It is designed to store information such as ID or object describing a resource during a test.
This metadata can use the `render` method to provide helpful golang templating features to have commands arguments computed dynamically.

#### BeforeFunc and AfterFunc

Usually, you might need to set up and teardown resources when you are running a test for testing a specific command.
For that you can use [`BeforeFunc`](https://github.com/scaleway/scaleway-cli/blob/master/internal/core/testing.go#L100) and [`AfterFunc`](https://github.com/scaleway/scaleway-cli/blob/master/internal/core/testing.go#L102).
Those types allow you to execute code before (`BeforeFunc`) or after (`AfterFunc`) the main command you want to test.
Those functions can access the metadata to change dynamically their behavior.

#### Logging and debug mode

When you are running CLI commands, you can use the `-D` to access the user logs.
You can activate the debug mode by passing the `-debug` flag when running the `go test` command.
For instance:

```
go test ./internal/namespaces/init -debug
```

When you are developing tests, you can also use the `Logger` field that is available in the different contexts to write your own logs.

### Checking the version

When you are running CLI commands the version check is made in every action. You can avoid these output setting `SCW_DISABLE_CHECK_VERSION` to false.

#### Targeting specific tests

The test suite is design to run quickly, but when you are recording interactions you probably don't want to record interactions for all the tests in the test suite.
You can filter the test you want to run using the native filtering features of the go test command.

So let's suppose you would like to run the test `Test_InstallServer` in the baremetal package, you would use:

`go test ./internal/namespaces/baremetal/v1 -run Test_InstallServer`

Keep in mind that running a single file is NOT equivalent to run the test for a package.
Always run the test on the whole package (here the "baremetal" package stored in the folder "./internal/namespaces/baremetal/v1") and use the `-run` to target specific tests.

### Adding new tests

We welcome contributions!
If you want to contribute new tests you should have the following:

1. Setup your dev environment:
    - Install go for your platform
    - Install your credentials, preferably in a configuration file (run `scw init`)
        - Keep in mind that if you record interaction, the resource you will instantiate will be delivered and billed.
        - Clean up the resource you don't use once the recording is over.