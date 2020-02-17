# CLI Release Process

This is the CLI v2 release process it will be scripted in `./scripts/release` in order to be automated.

The following process assumes that the user has 2 origins:

- `upstream`: `scaleway/scaleway-cli`
- `origin`: `<user>/scaleway-cli`

## STEP 1: Prepare the release

1. `cd` to git repository
2. `git checkout v2`
3. `git fetch upstream v2 && git rebase upstream/v2`
4. `git checkout -b release`
5. Update all occurences of the semver version in the README (with `sed`)
6. Update `Version` variable in `main.go`: increment version and remove `+dev` suffix
7. Update `CHANGELOG.md`
   1. Replace the `(Unreleased)` by the release date
   2. Fill the section with changes (template TDB)
8. `git commit -am "chore: release {semver}"`
9. `git push origin release`
10. Open a PR targeted to `v2` branch

## STEP 2: Make the Release

1. Merge the PR with Github
2. `git checkout v2`
3. `git fetch upstream v2 && git rebase upstream/v2`
4. `git tag {semver-version}`
5. `git push --tags upstream`
6. Build artifacts `./scripts/build.sh`
7. Create a Github release
   1. Name `{semver-version}`
   2. Tag:  `{semver-version}`
   3. Content: the content of the new `CHANGELOG.md` entry
   4. Check "this is a pre-release" if server contains `alpha`, `beta` or `rc`
   5. Add artifacts (compiled binaries + `SHA256SUMS`)
   6. Click on `Publish Release` button

## STEP 3: Cleanup the Release

1. `git checkout v2`
2. `git fetch upstream v2 && git rebase upstream/v2`
3. `git checkout -b post-release`
4. Add a `## {previous semver}+dev (Unreleased)` heading in `CHANGELOG.md`
5. Append `+dev` suffix in `Version` in `main.go`
6. Fix [outdated check version golden](https://github.com/scaleway/scaleway-cli/blob/v2/internal/core/testdata/test-check-version-outdated-version.stderr.golden)
7. `git commit -am "chore: cleanup after {semver-version}"`
8. `git push origin post-release`
9. Upgrade S3 version file

