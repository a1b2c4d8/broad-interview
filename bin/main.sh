#!/usr/bin/env bash

_BIN_DIR="$( cd "$( dirname "$0" )" && pwd )"
_ROOT_DIR="$( cd "$_BIN_DIR"/.. && pwd )"
_LOCAL_CMD="$_ROOT_DIR"/_output/main

_IMAGE_NAME="a1b2c4d8/broad-interview"
_GINKGO_TEST_SUITE="internal/service"

_IMAGE_NAME="a1b2c4d8/broad-interview"

_run() {
  cd "$_ROOT_DIR"
  if [[ "$1" == "--test" ]] ; then
    ( _test )
  else
    ( _main "$@" )
  fi

}

_ensure_docker_image() {
  if ! which docker > /dev/null ; then
    echo "Docker not found!"
    exit 1
  fi

  if ! docker images | grep "$_IMAGE_NAME" > /dev/null ; then
    docker build . -t "$_IMAGE_NAME":latest || exit 1
    echo
  fi
}

_main() {
  _ensure_docker_image
  _CMD="$( grep '^CMD' Dockerfile | awk '{ print $3 }' | sed 's/"//g' )"
  docker run --rm -it "$_IMAGE_NAME":latest "$_CMD" "$@"
}

_test() {
  _ensure_docker_image
  docker run --rm -it "$_IMAGE_NAME":latest ginkgo "$_GINKGO_TEST_SUITE"
}

( _run "$@" )
