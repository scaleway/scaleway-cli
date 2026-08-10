---
name: nightly-test-debugger
description: Analyze and debug failing nightly acceptance tests in this Go CLI repository. Use this skill whenever the user mentions nightly tests failing, CI failures in nightly.yml, acceptance test failures, or wants to investigate why a test failed in the nightly CI run. This includes diagnosing test failures, understanding error logs, identifying root causes, and suggesting fixes.
---

# Nightly Test Debugger

This skill helps you systematically analyze and debug failing nightly acceptance tests in the Scaleway CLI repository.

## Understanding Nightly Tests

The nightly test suite (`.github/workflows/nightly.yml`) runs acceptance tests for all products daily at 00:00 UTC. Each product is tested in isolation with:

- `CLI_UPDATE_CASSETTES: true` - Records API interactions
- `CLI_UPDATE_GOLDENS: true` - Updates expected output files
- Real API calls against Scaleway infrastructure

Common failure modes:
1. **API changes** - Upstream API changed, test expectations no longer match
2. **Cassette mismatches** - Recorded HTTP interactions don't match current behavior
3. **Golden file drift** - Output format changed
4. **Infrastructure issues** - Resource provisioning failures, quotas, timeouts
5. **Race conditions** - Async operations not completing in expected time
6. **Dependency changes** - Go module updates breaking compatibility

## Debugging Workflow

### Step 1: Gather Context

When investigating a failing nightly test, first collect the ID of the most recent run using:

gh run list --workflow "Nightly Acceptance Tests" --limit 1

### Step 2: Read the Failure Logs

Fetch the full logs from the GitHub Actions run:

```bash
gh run view <RUN_ID> --log | grep FAIL:
```

This will give you names of the failed tests.

### Step 3: Locate the Test File

Test files live in:
```text
internal/namespaces/<product>/v<version>/<test_file>.go
```

Common patterns:
- `custom_<feature>_test.go` - Feature-specific tests
- `<product>_test.go` - General tests

Read the failing test function to understand:
- What the test is trying to do
- What `BeforeFunc` setup runs
- What `Check` assertions are made
- What `AfterFunc` cleanup runs

### Step 4: Categorize the Failure

#### Timeout/Infrastructure Failure

**Symptoms:**
```text
context deadline exceeded
timeout waiting for resource
```

**Diagnosis:**
1. Check if the timeout value is reasonable
2. Look for resource provisioning delays
3. Check for quota issues in the nightly environment

**Fix:**
- Increase timeout in test config
- Add retry logic in `BeforeFunc`
- Check nightly credentials have proper quotas

## Resources

- Nightly workflow: `.github/workflows/nightly.yml`
- Test framework: `core/test_*.go`
- Example tests: `internal/namespaces/*/v*/*_test.go`
- Testdata: `internal/namespaces/*/v*/testdata/`

## When to Escalate

Some failures may not be fixable in the test itself:

1. **Upstream API breaking change** - May need SDK update
2. **Service degradation** - Contact Scaleway platform team
3. **Credential issues** - Check nightly secrets in GitHub

In these cases, clearly document:
- What you investigated
- What you ruled out
- What needs external action
