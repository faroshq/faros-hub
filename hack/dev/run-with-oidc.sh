#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace


source .env
go run ./cmd/hub-api start -a \
    --oidc-issuer-url=$FAROS_OIDC_ISSUER_URL \
    --oidc-client-id=$FAROS_OIDC_CLIENT_ID \
    --oidc-ca-file=$FAROS_OIDC_CA_FILE \
    --oidc-issuer-url=$FAROS_OIDC_ISSUER_URL \
    --oidc-username-claim=$FAROS_OIDC_USERNAME_CLAIM \
    --oidc-groups-claim=$FAROS_ODIC_GROUPS_CLAIM \
    "--oidc-username-prefix=$FAROS_OIDC_USER_PREFIX:" \
    "--oidc-groups-prefix=$FAROS_OIDC_GROUPS_PREFIX:"
