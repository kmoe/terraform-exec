#!/bin/bash

set -e
set -x

# release.sh will:
# 1. Modify changelog
# 2. Run changelog links script
# 3. Modify version in internal/version/version.go
# 4. Commit and push changes
# 5. Create a Git tag

function pleaseUseGNUsed {
    echo "Please install GNU sed to your PATH as 'sed'."
    exit 1
}

function init {
  sed --version > /dev/null || pleaseUseGNUsed

  DATE=$(date '+%B %d, %Y')

  if [ "$CI" = true ] ; then
    GPG_KEY_ID=C6DC8F8C8E78B36A
    gpg --batch --import <(echo -e "${GPG_PUBLIC_KEY}")
    gpg --batch --import <(echo -e "${GPG_PRIVATE_KEY}")
    git config --global user.email hashibot-feedback+tf-sdk-circleci@hashicorp.com
    git config --global user.name "Terraform SDK CircleCI"
  fi

  TARGET_VERSION="$(getTargetVersion)"
}

semverRegex='\([0-9]\+\.[0-9]\+\.[0-9]\+\)\(-\?\)\([0-9a-zA-Z.]\+\)\?'

function getTargetVersion {
  # parse target version from CHANGELOG
  sed -n 's/^# '"$semverRegex"' (Unreleased)$/\1\2\3/p' CHANGELOG.md || \
     (printf "\nTarget version not found in changelog, exiting" && \
       exit 1)
}

function modifyChangelog {
  sed -i "s/$TARGET_VERSION (Unreleased)$/$TARGET_VERSION ($DATE)/" CHANGELOG.md
}

function changelogLinks {
  ./scripts/release/changelog_links.sh
}

function changelogMain {
  printf "Modifying Changelog..."
  modifyChangelog
  printf "ok!\n"
  printf "Running Changelog Links..."
  changelogLinks
  printf "ok!\n"
}

function modifyVersionFiles {
  sed -i "s/const version =.*/const version = \"${TARGET_VERSION}\"/" internal/version/version.go
}

function commitChanges {
  git add CHANGELOG.md
  modifyVersionFiles
  git add internal/version/version.go

  if [ "$CI" = true ] ; then
      git commit --gpg-sign="${GPG_KEY_ID}" -m "v${TARGET_VERSION} [skip ci]"
      git tag -a -m "v${TARGET_VERSION}" -s -u "${GPG_KEY_ID}" "v${TARGET_VERSION}"
  else
      git commit -m "v${TARGET_VERSION} [skip ci]"
      git tag -a -m "v${TARGET_VERSION}" -s "v${TARGET_VERSION}"
  fi

  git push origin "${CIRCLE_BRANCH}"
  git push origin "v${TARGET_VERSION}"
}

function commitMain {
  printf "Committing Changes..."
  commitChanges
  printf "ok!\n"
}

function createGitHubRelease {
    curl -sL https://git.io/goreleaser | bash
}

function main {
  init
  changelogMain
  commitMain
  createGitHubRelease
}

main
