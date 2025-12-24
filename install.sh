#!/bin/bash

set -e

REPO="rakshit-gen/openskill"
SITE_URL="https://openskill.online"
BINARY_NAME="openskill"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo "  ╭─────────────────────────────────────╮"
echo "  │       OpenSkill CLI Installer       │"
echo "  │   Skill Management for Claude Code  │"
echo "  ╰─────────────────────────────────────╯"
echo ""

# Check for required dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo -e "${RED}Error: curl is required but not installed${NC}" >&2
    exit 1
fi

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Darwin) os="darwin" ;;
        Linux) os="linux" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *) echo -e "${RED}Unsupported OS: $(uname -s)${NC}" >&2; exit 1 ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) echo -e "${RED}Unsupported architecture: $(uname -m)${NC}" >&2; exit 1 ;;
    esac

    echo "${os}_${arch}"
}

# Get latest release version
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
}

PLATFORM=$(detect_platform)
echo -e "  Platform: ${GREEN}${PLATFORM}${NC}"

echo -e "  Fetching latest version..."
VERSION=$(get_latest_version)

if [ -z "$VERSION" ]; then
    echo -e "${YELLOW}  No releases found. Building from source...${NC}"

    # Check for Go
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}Error: Go is required to build from source${NC}" >&2
        echo "  Install Go from https://go.dev/dl/"
        exit 1
    fi

    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    echo "  Cloning repository..."
    git clone --depth 1 "https://github.com/${REPO}.git" openskill
    cd openskill/OpenSkill-cli

    echo "  Building..."
    go build -o "$BINARY_NAME" ./cmd/openskill

    echo "  Installing to ${INSTALL_DIR}..."
    sudo mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"

    cd /
    rm -rf "$TEMP_DIR"
else
    echo -e "  Version: ${GREEN}${VERSION}${NC}"

    # Construct download URL
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}_${PLATFORM}"

    if [[ "$PLATFORM" == windows* ]]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
    fi

    echo "  Downloading..."
    TEMP_FILE=$(mktemp)

    if ! curl -fsSL -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
        echo -e "${YELLOW}  Binary not found. Building from source...${NC}"

        if ! command -v go >/dev/null 2>&1; then
            echo -e "${RED}Error: Go is required to build from source${NC}" >&2
            exit 1
        fi

        TEMP_DIR=$(mktemp -d)
        cd "$TEMP_DIR"
        git clone --depth 1 --branch "$VERSION" "https://github.com/${REPO}.git" openskill 2>/dev/null || \
        git clone --depth 1 "https://github.com/${REPO}.git" openskill
        cd openskill/OpenSkill-cli
        go build -o "$BINARY_NAME" ./cmd/openskill
        sudo mv "$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
        cd /
        rm -rf "$TEMP_DIR"
    else
        chmod +x "$TEMP_FILE"
        echo "  Installing to ${INSTALL_DIR}..."
        sudo mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
fi

# Verify installation
if command -v openskill >/dev/null 2>&1; then
    echo ""
    echo -e "${GREEN}  ✓ OpenSkill installed successfully!${NC}"
    echo ""
    echo "  Get started:"
    echo "    openskill init          # Initialize in your project"
    echo "    openskill add my-skill  # Create a new skill"
    echo "    openskill list          # View all skills"
    echo ""
else
    echo -e "${RED}  Installation may have failed. Check ${INSTALL_DIR}${NC}" >&2
    exit 1
fi
