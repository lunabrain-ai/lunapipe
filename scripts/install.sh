#!/usr/bin/env bash

# Some helpful functions
yell() { echo -e "${RED}FAILED> $* ${NC}" >&2; }
die() { yell "$*"; exit 1; }
try() { "$@" || die "failed executing: $*"; }
log() { echo -e "--> $*"; }

# Colors for colorizing
RED='\033[0;31m'
GREEN='\033[0;32m'
PURPLE='\033[0;35m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'

INSTALL_PATH=${INSTALL_PATH:-"/usr/local/bin"}
NEED_SUDO=0

function maybe_sudo() {
    if [[ "$NEED_SUDO" == '1' ]]; then
        sudo "$@"
    else
        "$@"
    fi
}

# check for curl
hasCurl=$(which curl 2>/dev/null)
if [ "$?" = "1" ]; then
    die "You need to install curl to use this script."
fi

log "Selecting version..."

owner="lunabrain-ai"
repo="lunapipe"

# Get the releases JSON from the API
releases=$(curl -s "https://api.github.com/repos/${owner}/${repo}/releases")

# Use jq to get the latest release tag name
latest=$(echo $releases | jq -r '.[0].tag_name')

version=$latest

if [ ! $version ]; then
    log "${YELLOW}"
    log "Failed while attempting to install $repo. Please manually install:"
    log ""
    log "1. Open your web browser and go to https://github.com/$owner/$repo/releases"
    log "2. Download the cli from latest release for your platform. Name it '$repo'."
    log "3. chmod +x ./$repo"
    log "4. mv ./$repo /usr/local/bin"
    log "${NC}"
    die "exiting..."
fi

log "Selected version: $version"

log "${YELLOW}"
log NOTE: Install a specific version of the CLI by using VERSION variable
log "curl -L https://raw.githubusercontent.com/$owner/$repo/$version/scripts/install.sh | VERSION=$version bash"
log "${NC}"

# check for existing installation
hasCli=$(which $repo 2>/dev/null)
if [ "$?" = "0" ]; then
    log ""
    log "${GREEN}You already have $repo at '${hasCli}'${NC}"
    export n=3
    log "${YELLOW}Downloading again in $n seconds... Press Ctrl+C to cancel.${NC}"
    log ""
    sleep $n
fi

# get platform and arch
platform='unknown'
unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
    platform='Linux'
elif [[ "$unamestr" == 'Darwin' ]]; then
    platform='Darwin'
fi

if [[ "$platform" == 'unknown' ]]; then
    die "Unknown OS platform"
fi

arch='unknown'
archstr=`uname -m`
if [[ "$archstr" == 'x86_64' ]]; then
    arch='x86_64'
elif [[ "$archstr" == 'arm64' ]] || [[ "$archstr" == 'aarch64' ]]; then
    arch='arm64'
else
    die "prebuilt binaries for $(arch) architecture not available, please try building from source https://github.com/${owner}/${repo}/"
fi

# some variables
suffix="${platform}_${arch}"
tmpDir=${mktemp -d}
if [ -e $tmpDir ]; then
    rm -rf $tmpDir
fi

targetFile=$tmpDir/$repo

log "${PURPLE}Downloading $repo for $platform_$arch to ${tmpDir}${NC}"
url=https://github.com/$owner/$repo/releases/download/$version/${repo}_${version}_${suffix}.tar.gz

try curl -sL "$url" | tar -xz -C $tmpDir
try chmod +x $targetFile

log "${GREEN}Download complete!${NC}"

# check for sudo
needSudo=$(mkdir -p ${INSTALL_PATH} && touch ${INSTALL_PATH}/.${repo}install &> /dev/null)
if [[ "$?" == "1" ]]; then
    NEED_SUDO=1
fi
rm ${INSTALL_PATH}/.${repo}install &> /dev/null

if [[ "$NEED_SUDO" == '1' ]]; then
    log
    log "${YELLOW}Path '$INSTALL_PATH' requires root access to write."
    log "${YELLOW}This script will attempt to execute the move command with sudo.${NC}"
    log "${YELLOW}Are you ok with that? (y/N)${NC}"
    read a
    if [[ $a == "Y" || $a == "y" || $a = "" ]]; then
        log
    else
        log
        log "  ${BLUE}sudo mv $tmpDir ${INSTALL_PATH}/${repo}${NC}"
        log
        die "Please move the binary manually using the command above."
    fi
fi

log "Moving cli from $tmpDir to ${INSTALL_PATH}"

try maybe_sudo mv $tmpDir ${INSTALL_PATH}/${repo}

log
log "${GREEN}hasura cli installed to ${INSTALL_PATH}${NC}"
log

if [ -e $tmpDir ]; then
    rm -rf $tmpDir
fi

sh -c "${repo} -h"

if ! $(echo "$PATH" | grep -q "$INSTALL_PATH"); then
    log
    log "${YELLOW}$INSTALL_PATH not found in \$PATH, you might need to add it${NC}"
    log
fi
