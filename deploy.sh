#!/usr/bin/env bash
# rizz release helper
#
# Usage:
#   ./deploy.sh patch       # 0.1.0 -> 0.1.1
#   ./deploy.sh minor       # 0.1.0 -> 0.2.0
#   ./deploy.sh major       # 0.1.0 -> 1.0.0
#   ./deploy.sh 1.2.3       # explicit version
#
# What it does:
#   1. Verifies the working tree is clean and on main
#   2. Pulls latest from origin
#   3. Runs go vet, go build, go test
#   4. Computes the next tag (or uses the one you passed)
#   5. Creates an annotated tag and pushes it
#   6. GitHub Actions picks up the tag push and runs GoReleaser

set -euo pipefail

RED=$'\033[1;31m'
GREEN=$'\033[1;32m'
YELLOW=$'\033[1;33m'
BLUE=$'\033[1;34m'
RESET=$'\033[0m'

info()  { printf '%s==>%s %s\n' "$BLUE" "$RESET" "$*"; }
ok()    { printf '%s%s%s %s\n' "$GREEN" "✓" "$RESET" "$*"; }
warn()  { printf '%s%s%s %s\n' "$YELLOW" "!" "$RESET" "$*" >&2; }
fail()  { printf '%s%s%s %s\n' "$RED" "✗" "$RESET" "$*" >&2; exit 1; }

usage() {
    cat <<'EOF'
usage: ./deploy.sh <patch|minor|major|x.y.z>

  patch   bump the last number       (0.1.0 -> 0.1.1)
  minor   bump the middle number     (0.1.0 -> 0.2.0)
  major   bump the first number      (0.1.0 -> 1.0.0)
  x.y.z   use exactly this version   (e.g. 1.2.3)
EOF
    exit 1
}

require_clean_tree() {
    if [ -n "$(git status --porcelain)" ]; then
        fail "working tree is not clean. commit or stash first."
    fi
    ok "working tree clean"
}

require_main_branch() {
    branch=$(git rev-parse --abbrev-ref HEAD)
    if [ "$branch" != "main" ]; then
        fail "not on main (on '$branch'). switch to main first."
    fi
    ok "on main"
}

pull_latest() {
    info "pulling latest from origin/main"
    git pull --ff-only origin main
    ok "up to date with origin/main"
}

run_checks() {
    info "go vet ./..."
    go vet ./...
    ok "go vet passed"

    info "go build ./..."
    go build ./...
    ok "go build passed"

    info "go test ./..."
    go test ./...
    ok "go test passed"
}

latest_tag() {
    git tag --list 'v*' --sort=-v:refname | head -n1
}

bump_version() {
    current=$1   # e.g. v0.1.0 or empty
    part=$2      # patch | minor | major

    if [ -z "$current" ]; then
        case "$part" in
            patch) echo "0.0.1" ;;
            minor) echo "0.1.0" ;;
            major) echo "1.0.0" ;;
        esac
        return
    fi

    stripped=${current#v}
    IFS='.' read -r major minor patch <<< "$stripped"

    case "$part" in
        patch) patch=$((patch + 1)) ;;
        minor) minor=$((minor + 1)); patch=0 ;;
        major) major=$((major + 1)); minor=0; patch=0 ;;
    esac

    echo "$major.$minor.$patch"
}

validate_explicit_version() {
    v=$1
    if ! [[ "$v" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        fail "invalid version: '$v' (expected x.y.z)"
    fi
}

tag_exists() {
    git rev-parse -q --verify "refs/tags/$1" >/dev/null 2>&1
}

confirm() {
    printf '%s' "$1 [y/N] "
    read -r answer
    case "$answer" in
        y|Y|yes|YES) return 0 ;;
        *) return 1 ;;
    esac
}

main() {
    [ $# -eq 1 ] || usage

    input=$1
    require_main_branch
    require_clean_tree
    pull_latest

    case "$input" in
        patch|minor|major)
            current=$(latest_tag || true)
            new_version=$(bump_version "$current" "$input")
            info "current tag: ${current:-none}"
            ;;
        *)
            validate_explicit_version "$input"
            new_version=$input
            ;;
    esac

    new_tag="v$new_version"

    if tag_exists "$new_tag"; then
        fail "tag $new_tag already exists"
    fi

    run_checks

    info "next tag: $new_tag"
    confirm "create and push $new_tag?" || fail "aborted"

    info "creating tag $new_tag"
    git tag -a "$new_tag" -m "release $new_tag"

    info "pushing tag to origin"
    git push origin "$new_tag"

    ok "tag $new_tag pushed"
    info "GitHub Actions will now build and publish the release"
    info "  watch: https://github.com/Gogoro/rizz/actions"
    info "  releases: https://github.com/Gogoro/rizz/releases"
}

main "$@"
