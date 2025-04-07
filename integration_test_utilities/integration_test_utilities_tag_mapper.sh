#!/bin/sh

check_and_install_yq() {
    if ! command -v yq >/dev/null 2>&1; then
        echo "yq is not installed. Installing yq..."
        if [ "$(uname)" = "Darwin" ]; then
            if command -v brew >/dev/null 2>&1; then
                brew install yq
            else
                echo "Homebrew is not installed. Please install Homebrew first: https://brew.sh/"
                exit 1
            fi
        elif [ "$(uname)" = "Linux" ]; then
            if command -v wget >/dev/null 2>&1; then
                ARCH=$(uname -m)
                if [ "$ARCH" = "x86_64" ]; then
                    BINARY_ARCH="amd64"
                elif [ "$ARCH" = "aarch64" ]; then
                    BINARY_ARCH="arm64"
                else
                    echo "Unsupported architecture: $ARCH"
                    exit 1
                fi
                sudo wget -qO /usr/local/bin/yq "https://github.com/mikefarah/yq/releases/latest/download/yq_linux_$BINARY_ARCH"
                sudo chmod a+x /usr/local/bin/yq
            else
                echo "wget is not installed. Please install wget first."
                exit 1
            fi
        else
            echo "Unsupported operating system"
            exit 1
        fi
    fi
}

check_and_install_yq


product_mappings=""


{ git diff main...$(git rev-parse --abbrev-ref HEAD) --name-only; git diff --name-only; } | sort -u | grep '^newrelic/' | while read -r file; do

    filename=$(echo "$file" | sed 's|^newrelic/||')
    mapping=$(yq eval ".\"$filename\".product_mapping" "integration_test_utilities/integration_test_mappings.yaml")

    if [ "$mapping" = "null" ] || [ -z "$mapping" ]; then
        continue
    fi

    if ! echo "$product_mappings" | grep "$mapping"; then
        if [ -z "$product_mappings" ]; then
            product_mappings="$mapping"
        else
            product_mappings="$product_mappings,$mapping"
        fi
    fi
if [ -z "$product_mappings" ]; then
    echo "NONE"
else
    echo "-tags=$product_mappings"
fi
done
