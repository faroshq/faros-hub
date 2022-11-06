#!/bin/bash

DIR=hack/dev/dex/ssl
mkdir -p $DIR

cat << EOF > $DIR/req.cnf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = dex.dev.faros.sh
EOF



openssl genrsa -out $DIR/ca-key.pem 2048
openssl req -x509 -new -nodes -key $DIR/ca-key.pem -days 365 -out $DIR/ca.pem -subj "/CN=kube-ca"

openssl genrsa -out $DIR/key.pem 2048
openssl req -new -key $DIR/key.pem -out $DIR/csr.pem -subj "/CN=kube-ca" -config $DIR/req.cnf
openssl x509 -req -in $DIR/csr.pem -CA $DIR/ca.pem -CAkey $DIR/ca-key.pem -CAcreateserial -out $DIR/cert.pem -days 365 -extensions v3_req -extfile $DIR/req.cnf
