#!/usr/bin/env bash
set -euo pipefail

BRANCH=$(git branch --show-current)
if [[ "$BRANCH" != "main" ]]; then
  echo "Current branch is $BRANCH. Release logic only runs on main."
  exit 0
fi

git fetch --tags
TAG=$(git tag | sort -uV | tail -n 1)
[[ -z "${TAG:-}" ]] && TAG="v0.0.0"
echo "Latest tag: $TAG"

if [[ "$TAG" == "v0.0.0" ]]; then
  LOG=$(git log --pretty=format:%s)
else
  LOG=$(git log "$TAG..HEAD" --pretty=format:%s)
fi

BREAK=$(echo "$LOG" | grep -e '!' || true)
FEAT=$(echo "$LOG" | grep -e '^feat' || true)
FIX=$(echo "$LOG" | grep -e '^fix' || true)

VERSION=${TAG#v}
IFS='.' read -r major minor patch <<< "$VERSION"

if [[ -n "$BREAK" ]]; then
  major=$((major+1)); minor=0; patch=0
elif [[ -n "$FEAT" ]]; then
  minor=$((minor+1)); patch=0
elif [[ -n "$FIX" ]]; then
  patch=$((patch+1))
else
  echo "No relevant changes (fix/feat/break). Skipping release."
  exit 0
fi

NEWTAG="v$major.$minor.$patch"
echo "New version: $NEWTAG"

if [[ "$NEWTAG" == "$TAG" ]]; then
  echo "Nothing to tag."
  exit 0
fi

if [[ -n "${CI:-}" ]]; then
  git config user.name "GitHub Actions"
  git config user.email "actions@github.com"
fi

git tag "$NEWTAG"
git push origin "$NEWTAG"

echo "Pushed tag $NEWTAG"
