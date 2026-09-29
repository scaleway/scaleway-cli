# Testing the CLI

TL;DR

| `mise run test:cli` flag | Description                              |
|--------------------------|------------------------------------------|
| `--debug`                | Enable debug mode for the test           |
| `--goldens`              | Update the golden files of the run tests |
| `--cassettes`            | Update the cassettes of the run tests    |
| `--run <regex>`          | Run only tests matching the given regex  |
| `--race`                 | Enable data race detection               |
| `--format <fmt>`         | gotestsum output format                  |
| `--timeout <duration>`   | Test timeout (e.g. `30s`, `5m`, `1h`)    |

All flags are optional and can be combined. Run `mise run test:cli --help` for the full list.

## Objectives of the test suite

- Ensure that we got no new regression when we merge new code
- Avoid leaking credentials when doing integration testing with Scaleway APIs

## Interaction Recording: Golden & Cassettes files

A **cassette** file contains all the interactions of the Scaleway APIs that a CLI must produce when a given test run.
A **golden** file contains the output that the CLI must produce when a given interaction recorded by a cassette occurs.

Having golden ensure that any change in a command output will be noticed if the behavior of the CLI change for a given interaction.
Having cassette ensure that we can replay interactions without actually doing the calls.

Warning, if you choose to record a new cassette, you will create real resources on your organization and will be billed accordingly.
Be sure to check out the resource you create in each test and remember to delete them once you need them anymore.

## Metadata

When running a test, you might need information such as ID that you cannot know in advance (such as ID of resources).
The `core.Test` uses different helpers to pass useful information around.

One of them is the [`core.testMetadata`](https://github.com/scaleway/scaleway-cli/blob/main/internal/core/testing.go#L80).
It is designed to store information such as ID or object describing a resource during a test.
This metadata can use the `render` method to provide helpful golang templating features to have commands arguments computed dynamically.

## BeforeFunc and AfterFunc

Usually, you might need to set up and teardown resources when you are running a test for testing a specific command.
For that you can use [`BeforeFunc`](https://github.com/scaleway/scaleway-cli/blob/main/internal/core/testing.go#L100) and [`AfterFunc`](https://github.com/scaleway/scaleway-cli/blob/main/internal/core/testing.go#L102).
Those types allow you to execute code before (`BeforeFunc`) or after (`AfterFunc`) the main command you want to test.
Those functions can access the metadata to change dynamically their behavior.

## Logging and debug mode

When you are running CLI commands, you can use the `-D` to access the user logs.
You can activate the debug mode by passing the `--debug` flag to the test task.
For instance:

```shell
mise run test:cli --debug ./internal/namespaces/init
```

When you are developing tests, you can also use the `Logger` field that is available in the different contexts to write your own logs.

## Checking the version

When you are running CLI commands the version check is made in every action. You can avoid these output setting `SCW_DISABLE_CHECK_VERSION` to false.

## Targeting specific tests

The test suite is design to run quickly, but when you are recording interactions you probably don't want to record interactions for all the tests in the test suite.
You can filter the test you want to run using the `--run` flag of the test task.

So let's suppose you would like to run the test `Test_InstallServer` in the baremetal package, you would use:

`mise run test:cli --run Test_InstallServer ./internal/namespaces/baremetal/v1`

Keep in mind that running a single file is NOT equivalent to run the test for a package.
Always run the test on the whole package (here the "baremetal" package stored in the folder "./internal/namespaces/baremetal/v1") and use the `--run` to target specific tests.

## Adding new tests

We welcome contributions!
If you want to contribute new tests you should have the following:

1. Setup your dev environment:
    - Install [mise](https://mise.jdx.dev/) and run `mise install` to fetch the required tools (Go, golangci-lint, gotestsum, etc.)
    - Run `mise tasks ls` to list all available tasks, and `mise run <task> --help` to see the flags for a specific task.
    - Install your credentials, preferably in a configuration file (run `scw init`)
        - Keep in mind that if you record interaction, the resource you will instantiate will be delivered and billed.
        - Clean up the resource you don't use once the recording is over.

## Recording IAM API keys test cassettes (`Test_iamAPIKeyGet`)

`Test_iamAPIKeyGet` in `internal/namespaces/iam/v1alpha1/custom_iam_test.go` checks the output of `scw iam
api-key get` for three kinds of principals: the organization owner, an IAM member user, and an IAM application.

- **`GetOwnerAPIKey`**: no setup needed, it resolves the owner and one of its API keys at recording time.
- **`GetApplicationAPIKey`**: fully automated. The application and its API key are created via the CLI in the
  `BeforeFunc` and deleted in the `AfterFunc` (only when recording cassettes), so nothing to prepare manually.
- **`GetMemberAPIKey`**: requires a one-time manual setup (see below). The test picks any member of your
  organization (with an API key) at recording time. If no suitable member is present, the subtest is **skipped**
  and a message points to this section.

### Why the member cannot be automated

Recording the member scenario needs an *activated* member user with an API key in your organization. This cannot
be created and cleaned up from the CLI or the API:

- IAM users can only be **invited** to an organization; activation (accepting the invitation, setting a password)
  happens in the Scaleway console and cannot be scripted.
- The organization owner cannot create an API key for another member (no impersonation).
- Members **cannot be deleted** through the IAM API (only guest users can).

There is therefore no dynamic (create/activate/delete) member bootstrap, regardless of e-mail delivery. The
member is a one-time setup, shared across all future cassette recordings of this test. The test does not rely on
a specific member ID: it uses the first member returned by `scw iam user list type=member` and skips when none
is present.

### One-time member setup

1. In the [Scaleway console](https://console.scaleway.com/iam), go to **IAM → Users** and invite a user with an
   e-mail address you control (e.g. `testiam@testiam.testiam`). Nothing is sent to that mailbox: activation is
   done from the console of the invited account.
2. Open the invitation and complete the activation so the user becomes an *activated* member (`status: activated`
   in `scw iam user list`).
3. Log in with the member account and create an API key for it (or create it from the console under **IAM →
   API keys** while connected as the member):
   `scw iam api-key create user-id=<member-user-id> description=test-cli-iam-membership-api-key`
4. (Optional) Attach a policy to the member so it is displayed in the `api-key get` output:
   `scw iam policy create name=test-cli-iam-member user-id=<member-user-id> rules.0.permission-set-names.0=IAMManager`

### Recording the cassettes

```shell
mise run test:cli --cassettes --run Test_iamAPIKeyGet ./internal/namespaces/iam/v1alpha1
mise run test:cli --goldens --run Test_iamAPIKeyGet ./internal/namespaces/iam/v1alpha1
```

The `GetApplicationAPIKey` subtest creates and removes its own resources during the recording; only the member
setup above is manual. Remember to clean up the member's API key and any attached policy once you are done
recording, as members themselves cannot be removed through the API.
