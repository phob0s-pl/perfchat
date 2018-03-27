#!/bin/bash

set -e

main() {
  if ! which p2psim &>/dev/null; then
    fail "missing p2psim binary (you need to build cmd/p2psim and put it in \$PATH)"
  fi

  info "creating ${1} nodes"
  for i in $(seq 1 ${1}); do
    p2psim node create --name "$(node_name $i)"
    p2psim node start "$(node_name $i)"
  done

  info "connecting nodes to all other nodes"
  for i in $(seq 1 ${1}); do
        for j in $(seq 1 ${1}); do
		if [ $i -ne $j ]; then
		    p2psim node connect "$(node_name $j)" "$(node_name $i)"
		fi
	done
  done

  info "done"
}

node_name() {
  local num=$1
  echo "node$(printf '%02d' $num)"
}

info() {
  echo -e "\033[1;32m---> $(date +%H:%M:%S) ${@}\033[0m"
}

fail() {
  echo -e "\033[1;31mERROR: ${@}\033[0m" >&2
  exit 1
}

main "$1"
