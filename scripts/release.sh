#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_NONE='\033[0m'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ $CURRENT_GIT_BRANCH != 'main' ]; then
  printf "\n"
  printf "${COLOR_RED} Error: The release.sh script must be run while on the main branch. \n ${COLOR_NONE}"
  printf "\n"

  exit 1
fi

if [ $# -ne 1 ]; then
  printf "\n"
  printf "${COLOR_RED} Error: Release version argument required. \n\n ${COLOR_NONE}"
  printf " Example: \n\n    ./scripts/release.sh 0.9.0 \n\n"
  printf "  Example (make): \n\n    make release version=0.9.0 \n"
  printf "\n"

  exit 1
fi

RELEASE_VERSION=$1
GIT_USER=$(git config user.email)

echo "Generating release for v${RELEASE_VERSION} using git user ${GIT_USER}"

git checkout -b release/v${RELEASE_VERSION}

# Auto-generate CHANGELOG updates
git-chglog --next-tag v${RELEASE_VERSION} -o CHANGELOG.md --sort semver
# Fix any spelling issues in the CHANGELOG
misspell -source text -w CHANGELOG.md

# Commit CHANGELOG updates
git add CHANGELOG.md
git commit --no-verify -m "chore(changelog): Update CHANGELOG for v${RELEASE_VERSION}"
git push --no-verify origin release/v${RELEASE_VERSION}
