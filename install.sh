#!/usr/bin/env sh
# rizz installer
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Gogoro/rizz/main/install.sh | sh
#
# Environment:
#   RIZZ_VERSION   pin a specific version (e.g. v0.1.0). default: latest
#   RIZZ_INSTALL   install directory. default: ~/.local/bin (falls back to /usr/local/bin with sudo)

set -eu

REPO="Gogoro/rizz"
BINARY="rizz"

info()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33m!!!\033[0m %s\n' "$*" >&2; }
fail()  { printf '\033[1;31mxxx\033[0m %s\n' "$*" >&2; exit 1; }

detect_os() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        darwin) echo "darwin" ;;
        linux)  echo "linux" ;;
        *)      fail "unsupported OS: $os (try: go install github.com/$REPO@latest)" ;;
    esac
}

detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64) echo "x86_64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) fail "unsupported arch: $arch" ;;
    esac
}

fetch_latest_version() {
    url="https://api.github.com/repos/$REPO/releases/latest"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" | grep '"tag_name"' | head -n1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$url" | grep '"tag_name"' | head -n1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
    else
        fail "need curl or wget"
    fi
}

download() {
    src=$1
    dest=$2
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$src" -o "$dest"
    else
        wget -q "$src" -O "$dest"
    fi
}

pick_install_dir() {
    if [ -n "${RIZZ_INSTALL:-}" ]; then
        echo "$RIZZ_INSTALL"
        return
    fi
    if [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        echo "$HOME/.local/bin"
        return
    fi
    echo "/usr/local/bin"
}

install_binary() {
    src=$1
    dest_dir=$2
    mkdir -p "$dest_dir" 2>/dev/null || true
    if [ -w "$dest_dir" ]; then
        install -m 0755 "$src" "$dest_dir/$BINARY"
    else
        info "writing to $dest_dir requires sudo"
        sudo install -m 0755 "$src" "$dest_dir/$BINARY"
    fi
}

check_path() {
    dir=$1
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) warn "$dir is not on your PATH. add this to your shell rc:"
           warn "  export PATH=\"$dir:\$PATH\"" ;;
    esac
}

main() {
    os=$(detect_os)
    arch=$(detect_arch)
    version="${RIZZ_VERSION:-$(fetch_latest_version)}"
    [ -n "$version" ] || fail "could not determine latest version"

    # strip leading 'v' for the asset filename (goreleaser .Version)
    version_num=${version#v}

    asset="${BINARY}_${version_num}_${os}_${arch}.tar.gz"
    url="https://github.com/$REPO/releases/download/$version/$asset"

    info "installing $BINARY $version for $os/$arch"
    info "downloading $url"

    tmp=$(mktemp -d)
    trap 'rm -rf "$tmp"' EXIT

    download "$url" "$tmp/$asset"
    tar -xzf "$tmp/$asset" -C "$tmp"

    dest_dir=$(pick_install_dir)
    install_binary "$tmp/$BINARY" "$dest_dir"

    info "installed $BINARY to $dest_dir/$BINARY"
    check_path "$dest_dir"
    info "run: $BINARY --version"
}

main "$@"
