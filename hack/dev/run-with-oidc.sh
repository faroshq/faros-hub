#!/bin/bash

source .env
go run ./cmd/hub-api \
    --oidc-issuer-url=$FAROS_OIDC_ISSUER_URL \
    --oidc-client-id=$FAROS_OIDC_CLIENT_ID \
    --oidc-groups-claim=$FAROS_ODIC_GROUPS_CLAIM \
    --oidc-username-claim=$FAROS_OIDC_USERNAME_CLAIM \
    --oidc-username-prefix=$FAROS_OIDC_USERNAME_PREFIX \
    --oidc-groups-prefix=$FAROS_OIDC_GROUP_PREFIX
