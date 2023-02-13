#!/bin/sh

set -eou pipefail

pnpm run start --help &

/usr/local/bin/go-hasura-auth serve --debug
