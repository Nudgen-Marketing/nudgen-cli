#!/bin/bash
set -euo pipefail

REPO="Nudgen-Marketing/nudgen-cli"
BIN_NAME="nudgen"
BIN_DIR="${NUDGEN_BIN_DIR:-$HOME/.local/bin}"
VERSION="${NUDGEN_VERSION:-}"

if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
  bold()  { printf '\033[1m%s\033[0m' "$1"; }
  green() { printf '\033[32m%s\033[0m' "$1"; }
  red()   { printf '\033[31m%s\033[0m' "$1"; }
  dim()   { printf '\033[2m%s\033[0m' "$1"; }
  emerald() { printf '\033[38;2;16;185;129m%s\033[0m' "$1"; }
else
  bold()  { printf '%s' "$1"; }
  green() { printf '%s' "$1"; }
  red()   { printf '%s' "$1"; }
  dim()   { printf '%s' "$1"; }
  emerald() { printf '%s' "$1"; }
fi

info()  { echo "  $(emerald "✓") $1"; }
step()  { echo "  $(bold "→") $1"; }
error() { echo "  $(red "✗ ERROR:") $1" >&2; exit 1; }

find_sha256_cmd() {
  if command -v sha256sum &>/dev/null; then
    echo "sha256sum"
  elif command -v shasum &>/dev/null; then
    echo "shasum -a 256"
  else
    error "No SHA256 tool found (need sha256sum or shasum)"
  fi
}

detect_platform() {
  local os arch

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    darwin) os="darwin" ;;
    linux) os="linux" ;;
    freebsd) os="freebsd" ;;
    openbsd) os="openbsd" ;;
    mingw*|msys*|cygwin*) os="windows" ;;
    *) error "Unsupported OS: $os" ;;
  esac

  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) error "Unsupported architecture: $arch" ;;
  esac

  echo "${os}_${arch}"
}

get_latest_version() {
  local url version
  url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest" 2>/dev/null) || true
  version="${url##*/}"
  version="${version#v}"
  if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # Fallback to trial v1.0.0 if no tags found yet (better for initial repo setup)
    if [[ "$version" == "latest" ]] || [[ -z "$version" ]]; then
        echo "1.0.0"
        return
    fi
    error "Could not determine latest version (resolved '${version:-<empty>}' from '${url:-<no URL>}')."
  fi
  echo "$version"
}

verify_checksums() {
  local version="$1"
  local tmp_dir="$2"
  local archive_name="$3"
  local base_url="https://github.com/${REPO}/releases/download/v${version}"
  
  # Check if checksums.txt exists
  if ! curl -fsSLI "${base_url}/checksums.txt" &>/dev/null; then
    info "Skipping checksum verification (no checksums.txt in release)"
    return 0
  fi

  step "Verifying checksums..."
  if ! curl -fsSL "${base_url}/checksums.txt" -o "${tmp_dir}/checksums.txt"; then
    error "Failed to download checksums.txt"
  fi

  local expected actual
  expected=$(awk -v f="$archive_name" '$2 == f || $2 == ("*" f) {print $1; exit}' "${tmp_dir}/checksums.txt")
  actual=$(cd "$tmp_dir" && $(find_sha256_cmd) "$archive_name" | awk '{print $1}')
  [[ -n "$expected" && "$expected" == "$actual" ]]  \
    || error "Checksum verification failed for $archive_name"

  info "Checksum verified"
}

download_binary() {
  local version="$1"
  local platform="$2"
  local tmp_dir="$3"
  local url archive_name ext

  if [[ "$platform" == windows_* ]]; then
    ext="zip"
  else
    ext="tar.gz"
  fi

  archive_name="${BIN_NAME}_${version}_${platform}.${ext}"
  url="https://github.com/${REPO}/releases/download/v${version}/${archive_name}"

  step "Downloading ${BIN_NAME} v${version} for ${platform}..."

  if ! curl -fsSL "$url" -o "${tmp_dir}/${archive_name}"; then
    error "Failed to download from $url. Please make sure the release v${version} exists."
  fi

  verify_checksums "$version" "$tmp_dir" "$archive_name"

  step "Extracting..."
  if [[ "$ext" == "zip" ]]; then
    unzip -q "${tmp_dir}/${archive_name}" -d "$tmp_dir"
  else
    tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"
  fi

  local target_bin="${BIN_NAME}"
  if [[ "$platform" == windows_* ]]; then
    target_bin="${BIN_NAME}.exe"
  fi

  if [[ ! -f "${tmp_dir}/${target_bin}" ]]; then
    # Try to find it if it was in a subfolder
    find_bin=$(find "$tmp_dir" -name "$target_bin" -type f | head -n 1)
    if [[ -n "$find_bin" ]]; then
        mv "$find_bin" "${tmp_dir}/${target_bin}"
    else
        error "Binary not found in archive"
    fi
  fi

  mkdir -p "$BIN_DIR"
  mv "${tmp_dir}/${target_bin}" "$BIN_DIR/"
  chmod +x "$BIN_DIR/$target_bin"

  info "Installed ${BIN_NAME} to $BIN_DIR/$target_bin"
}

setup_path() {
  if [[ ":$PATH:" == *":$BIN_DIR:"* ]]; then
    return 0
  fi

  step "Adding $BIN_DIR to PATH"

  local shell_rc=""
  case "${SHELL:-}" in
    */zsh)  shell_rc="$HOME/.zshrc" ;;
    */bash) shell_rc="$HOME/.bashrc" ;;
    *)      shell_rc="$HOME/.profile" ;;
  esac

  local path_line="export PATH=\"$BIN_DIR:\$PATH\""

  if [[ -f "$shell_rc" ]] && grep -qF "$BIN_DIR" "$shell_rc" 2>/dev/null; then
    info "PATH already configured in $shell_rc"
  else
    echo "" >> "$shell_rc"
    echo "# Added by ${BIN_NAME} installer" >> "$shell_rc"
    echo "$path_line" >> "$shell_rc"
    info "Added to $shell_rc"
    info "Run: source $shell_rc"
  fi
}

verify_install() {
  local platform="$1"
  local target_bin="${BIN_NAME}"
  if [[ "$platform" == windows_* ]]; then
    target_bin="${BIN_NAME}.exe"
  fi

  local installed_version
  if installed_version=$("$BIN_DIR/$target_bin" version 2>/dev/null | head -n 1); then
    info "$(green "Nudgen CLI installed successfully!")"
    return 0
  fi

  error "Installation failed - binary not working"
}

show_banner() {
  local cols
  cols=$(tput cols 2>/dev/null || echo 80)
  if [[ "$cols" -ge 44 ]]; then
    local e="" b="" r=""
    if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
      e=$'\033[38;2;16;185;129m'  # Emerald 500
      b=$'\033[1m'                # bold
      r=$'\033[0m'                # reset
    fi

    local logo=(
      "   _  _ _  _ ___  ____ ____ _  _ "
      "   |\\ | |  | |  \\ | __ |___ |\\ | "
      "   | \\| |__| |__/ |__] |___ | \\| "
      "                                 "
    )

    echo ""
    if [[ -t 1 ]] && [[ -z "${NO_COLOR:-}" ]]; then
      for line in "${logo[@]}"; do
        echo "${e}${line}${r}"
        sleep 0.03
      done
    else
      for line in "${logo[@]}"; do
        echo "$line"
      done
    fi
    echo ""
  else
    echo ""
    echo "Nudgen CLI"
    echo ""
  fi
}

main() {
  show_banner

  if ! command -v curl &>/dev/null; then
    error "curl is required but not installed"
  fi

  local platform version tmp_dir
  platform=$(detect_platform)

  if [[ -n "$VERSION" ]]; then
    version="$VERSION"
  else
    version=$(get_latest_version)
  fi

  tmp_dir=$(mktemp -d)
  trap "rm -rf '${tmp_dir}'" EXIT

  download_binary "$version" "$platform" "$tmp_dir"
  setup_path
  verify_install "$platform"

  echo ""
  echo "  Next steps:"
  echo "    $(bold "nudgen login")        Authenticate with Nudgen"
  echo "    $(bold "nudgen --help")       Explore available commands"
  echo ""
}

main "$@"
