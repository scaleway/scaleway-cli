**This guide is intended for scaleway-cli maintainers. If you are not a maintainer, you probably want to check out the [documentation](README.md) instead.**

## Package release HOWTO

Ready to deploy a new version to users? Let's make a checklist of what you need to do.

For the sake of the example, we assume you want to release the version `42.8` and the previous version was `42.7`.

### Commit release

* Edit the Changelog at the bottom of [README.md](README.md).

  - Header: change the version from *v42.7* version) to *v42.8*
  - Header: replace "unreleased" with the current date
  - Content: read commits history, and fill the changelog if it's not already done
  - In the "To install a release" section, grep *v42.7* and replace with *v42.8*
  - Footer: edit the "View full commit list" link from */compare/v42.7..master* to */compare/v42.7..v42.8*

* Edit [pkg/scwversion/version.go](pkg/scwversion/version.go) and update the *VERSION* to *v42.8*.
* Make the commit release: `git commit -a -m 'Release v42.8'`.
* Tag the commit: `git tag v42.8`.
* Push: `git push && git push --tags`.

### Make a github release

* [Draft a new release](https://github.com/scaleway/scaleway-cli/releases) on Github.
* Build cross-compiled files: `make prepare-release-dist`. *This step relies on the go compiler, which is able to cross compile go projects. If you hack on the scaleway-cli in a Docker container, you can run this step from your container.*
* Build debian packages: `make prepare-release-debian-packages VERSION=42.8`. **The VERSION should not include the leading "v"**. *This step runs Docker containers. If you are hacking on scaleway-cli with Docker, you need to build the Debian packages from the Docker host, and not from your container.*
* Upload the Debian packages generated in `dist/latest`.
* Publish the release.

### Docker image

* Call `make prepare-release-docker-image VERSION=42.8`. **The VERSION should not include the leading "v"**. *This step runs Docker containers. If you are hacking on scaleway-cli with Docker, run this step from the Docker host and not from your container.*
* Push the local Docker image to the Docker hub: `docker push scaleway/cli`.

### Homebrew (OSX) package

* Get the released archive's sha256sum: `curl -sL https://github.com/scaleway/scaleway-cli/archive/v42.8.tar.gz | shasum -a 256`.
* Clone the homebrew Github repository: `git clone git@github.com:Homebrew/homebrew-core.git` to you personal account.
* Edit *Formula/scw.rb* and fix the *URL* **AND** the *sha256* **on top** of the file. You don't need to edit the SHAsums below. They will be updated automatically by Homebrew maintainers when the PR will be merged.
* Ensure the formula works: `brew install --build-from-source Formula/scw.rb`. You will probably need to uninstall your current installation of scaleway-cli before installing the formula.
* Commit with `git commit -a -m 'scw 42.8'` and make a pull request from your repository to [homebrew](https://github.com/Homebrew/homebrew-core) to make your new version official.

### Archlinux package

**This section is incomplete. Edit this part if you have additional informations.**

There is a Archlinux community package (aka. "AUR" — Archlinux User Repository) for scaleway-cli: https://aur.archlinux.org/packages/scaleway-cli/ maintained by "moscar". We should probably ping him when we make a new release.

### Update VERSION file

From time to time, scaleway-cli makes a HTTP query to `https://scw-devtools.s3.nl-ams.scw.cloud/scw-cli-version` to check if it is at the latest version available. This file needs to be updated.

**PLEASE WAIT UNTIL THE HOMEBREW PULL REQUEST IS ACCEPTED BEFORE UPDATING THE VERSION FILE.**


```
$> echo '42.8' > VERSION
$> s3cmd put VERSION s3://scaleway/scaleway-cli/VERSION --acl-public
# Ensure it's at the latest version
$> curl https://scw-devtools.s3.nl-ams.scw.cloud/scw-cli-version
# Should display "42.8".
```

### Post release commit

* Edit [README.md](README.md) and create a new *(unreleased)* entry:

   ```
   ### v42.8+dev (unreleased)

   This is the current development version. Update below with your changes. Remove
   this line when releasing the package.

   View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v42.8...master)
   ```

* Edit [pkg/scwversion/version.go](pkg/scwversion/version.go) and set *VERSION* to *v42.8+dev*.
