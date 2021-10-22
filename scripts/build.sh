#!/usr/bin/env bash
PKG=$5
BINARY_NAME=$2
function _info() {
  local msg=$1
  echo "[INFO] ${now} ${msg}"
}
function _version() {
  local msg=$1
  echo "[INFO] ${now} ${msg}"
}
function get_tag() {
  if [ -d ".git" ]; then
    local tag=$(git describe --tags)
    if ! [ $? -eq 0 ]; then
      local tag='unknown'
    else
      local tag=$(echo ${tag} | cut -d '-' -f 1)
    fi
    echo ${tag}
  fi
}
function get_branch() {
  if [ -d ".git" ]; then
    local branch=$(git rev-parse --abbrev-ref HEAD)
    if ! [ $? -eq 0 ]; then
      local branch='unknown'
    fi
    echo ${branch}
  fi
}
function get_commit() {
  if [ -d ".git" ]; then
    local commit=$(git rev-parse HEAD)
    if ! [ $? -eq 0 ]; then
      local commit='unknown'
    fi
    echo ${commit}
  fi
}
function build_in_docker() {
  docker run --rm -e 'GOOS=linux' -e 'GOARCH=amd64' \
    -v "$PWD":/go/src/${PKG} \
    -w /go/src/${PKG} golang:1.12.9 \
    sh -c "/bin/bash build/prepare.sh && go build -a -o ${bin_name} -ldflags \"-s -w\" -ldflags \"-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'\" ${main_file}"
}
function build() {
  local platform=$1
  local bin_name=$2
  local main_file=$3
  local image_prefix=$4
  local version=$(go version | grep -o 'go[0-9].[0-9].*')
  if [ ${platform} == "local" ]; then
    _info "Start local build ..."
    go build -o ${bin_name} -ldflags "-s -w" -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    _info "Program construction is complete: $2"
  elif [ ${platform} == "linux" ]; then
    _info "Start building a Linux platform version ..."
    GOOS=linux GOARCH=amd64 \
      go build -a -o ${bin_name} -ldflags "-s -w" -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    _info "Program construction is complete: $2"
  elif [ ${platform} == "docker" ]; then
    _info "Start based on Docker ..."
    build_in_docker
    _info "Program construction is complete: $2"
  elif [ ${platform} == "image" ]; then
    _info "Start building a Docker image ..."
    docker build . -t ${image_prefix}/${bin_name}:${TAG}
    _info "Clear intermediate image ..."
    docker ps -a | grep "Exited" | awk '{print $1 }' | xargs docker stop
    docker ps -a | grep "Exited" | awk '{print $1 }' | xargs docker rm
    docker rmi $(docker images -qf dangling=true) &>/dev/null
    _info "The Docker image is built: ${image_prefix}/${bin_name}:${TAG}"
  else
    echo "Please make sure the position variable is local, docker or linux."
  fi
}
function main() {
  _info "Start building [$2] ..."
  TAG=$(get_tag)
  BRANCH=$(get_branch)
  COMMIT=$(get_commit)
  DATE=$(date '+%Y-%m-%d %H:%M:%S')
  Path="${PKG}/version"
  _version "Build time (Build Time): $DATE"
  _version "Currently built version (Git Tag): $TAG"
  _version "Currently built branch (Git Branch): $BRANCH"
  _version "Commit of the current build (Git Commit): $COMMIT"
  build $1 $2 $3 $4
}
main $1 $2 $3 $4 $5
