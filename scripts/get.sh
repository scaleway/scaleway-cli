#!/bin/sh -e

CHOOSE_ANOTHER_INSTALLATION_METHOD="Please choose another installation method https://github.com/scaleway/scaleway-cli#installation."

# Detect the operating system and CPU architecture
echo "Detecting your operating system and CPU architecture..." >&2

case "$(uname -s)" in
Darwin) os="darwin" ;;
FreeBSD) os="freebsd" ;;
Linux) os="linux" ;;
*)
    echo "$(uname -s) is not supported by this installation script. $CHOOSE_ANOTHER_INSTALLATION_METHOD" >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
x86_64) arch="amd64" ;;
armv8*) arch="arm64" ;;
aarch64) arch="arm64" ;;
armv*) arch="arm" ;;
i386) arch="386" ;;
*)
    echo "$(uname -m) is not supported by this installation script. $CHOOSE_ANOTHER_INSTALLATION_METHOD" >&2
    exit 1
    ;;
esac

# Check if curl or wget is available
echo "Checking if curl or wget is available..." >&2

has_curl=$(which curl || true)
has_wget=$(which wget || true)

if [ -z "$has_curl" ] && [ -z "$has_wget" ]; then
    echo "You need curl or wget to proceed." >&2
    exit 1
fi

# Get the latest release from GitHub API
echo "Fetching the latest release from GitHub..." >&2

if [ -n "$has_curl" ]; then
    latest_release_json=$(curl -s https://api.github.com/repos/scaleway/scaleway-cli/releases/latest)
elif [ -n "$has_wget" ]; then
    latest_release_json=$(wget -q -O - https://api.github.com/repos/scaleway/scaleway-cli/releases/latest)
fi

latest=$(echo "$latest_release_json" | grep "browser_download_url.*${os}_${arch}" | cut -d : -f 2,3 | tr -d \" | tr -d " ")
if [ -z "$latest" ]; then
    echo "Unable to find the latest ${os}_${arch} release. https://github.com/scaleway/scaleway-cli/releases" >&2
    exit 1
fi

# Check if sudo or root permissions are available
echo "Checking if sudo or root permissions are available..." >&2

is_root=$(id -u)
has_sudo=$(which sudo || true)

if [ "$is_root" -eq 0 ]; then
    sudo=""
elif [ -n "$has_sudo" ]; then
    sudo="sudo"
    echo "The installation script will use sudo to proceed. You may be asked to enter your password." >&2
else
    echo "You need sudo or root permissions to proceed." >&2
    exit 1
fi

# Download the latest release
echo "Downloading the latest release..." >&2

if [ -n "$has_curl" ]; then
    $sudo curl -s -L "$latest" -o /usr/local/bin/scw
elif [ -n "$has_wget" ]; then
    $sudo wget -q -O /usr/local/bin/scw "$latest"
fi

echo "Setting executable permissions..." >&2
$sudo chmod +x /usr/local/bin/scw

# Success
echo "\033[0;32mScaleway CLI has been successfully installed! Start using it with: \033[1;32mscw --help\033[0m"
