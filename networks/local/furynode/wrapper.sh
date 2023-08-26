#!/usr/bin/env sh

BINARY=/fud/linux/${BINARY:-fud}
echo "binary: ${BINARY}"
ID=${ID:-0}
LOG=${LOG:-fud.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'fud' E.g.: -e BINARY=fud_my_test_version"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"

if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export FUDHOME="/fud/node${ID}/fud"

if [ -d "$(dirname "${FUDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${FUDHOME}" "$@" | tee "${FUDHOME}/${LOG}"
else
  "${BINARY}" --home "${FUDHOME}" "$@"
fi
