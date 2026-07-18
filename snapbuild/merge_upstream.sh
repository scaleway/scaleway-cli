#!/bin/bash
set -euo pipefail

UPSTREAM=upstream
UPSTREAM_URL=https://github.com/scaleway/scaleway-cli.git
BRANCH=main

if (( $# != 1 )); then
    >&2 echo "Syntax: $0 <tag>"
    exit 1
fi

TAG=$1

# Snap versions must not start with 'v' and may only contain [a-zA-Z0-9.+~-]
# Strip a leading 'v' from the tag (e.g. v2.35.0 -> 2.35.0).
SNAP_VERSION="${TAG#v}"

for cmd in git sed; do
    if ! command -v "$cmd" &> /dev/null; then
        >&2 echo "Required command '$cmd' not found"
        exit 1
    fi
done

# Ensure upstream remote exists
if ! git remote | grep -qx "$UPSTREAM"; then
    echo "Adding remote '$UPSTREAM' -> $UPSTREAM_URL"
    git remote add "$UPSTREAM" "$UPSTREAM_URL"
fi

echo "Fetching $UPSTREAM ..."
git fetch "$UPSTREAM" --tags

# Verify the tag exists on upstream
if ! git ls-remote --tags "$UPSTREAM" | grep -q "refs/tags/${TAG}$"; then
    >&2 echo "Tag '${TAG}' not found on upstream remote"
    exit 1
fi

echo "Checking out $BRANCH"
git checkout "$BRANCH"

echo "Pulling $BRANCH"
git pull origin "$BRANCH"

echo "Merging upstream tag ${TAG} into ${BRANCH}"
git merge "${TAG}" -m "Merge upstream tag ${TAG}"

# Re-apply our snapcraft files on top (merge may have overwritten them
# if upstream ever adds their own snapcraft.yaml).
echo "Generating snapcraft.yaml for version ${SNAP_VERSION}"
sed "s/VERSION/${SNAP_VERSION}/g" snapbuild/snapcraft.yaml.template > snapcraft.yaml

git add snapcraft.yaml snapbuild/snapcraft.yaml.template
git commit -m "Set snapcraft version to ${SNAP_VERSION} (upstream tag ${TAG})"

echo "Pushing to origin"
git push origin "$BRANCH"

echo "Done. Snap version will be: ${SNAP_VERSION}"
