#!/bin/bash
# Inspired from: https://github.com/grpc/grpc-java/tree/master/examples#generating-self-signed-certificates-for-use-with-grpc
# and https://github.com/grpc/grpc-java/tree/master/examples/example-tls

# Output files
# ca.key: Certificate Authority private key file (this shouldn't be shared in real-life)
# ca.crt: Certificate Authority trust certificate (this should be shared with users in real-life)
# server.key: Server private key, password protected (this shouldn't be shared)
# server.csr: Server certificate signing request (this should be shared with the CA owner)
# server.crt: Server certificate signed by the CA (this would be sent back by the CA owner) - keep on server
# server.pem: Conversion of server.key into a format gRPC likes (this shouldn't be shared)

# Summary 
# Private files: ca.key, server.key, server.pem, server.crt
# "Share" files: ca.crt (needed by the client), server.csr (needed by the CA)

# Change this CN to match your host in your environment if needed.
SERVER_CN="localhost"

# Create a configuration file for the extensions
cat <<EOF > server_cert_ext.cnf
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
x509_extensions = v3_req
[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default = US
stateOrProvinceName = State or Province Name (full name)
stateOrProvinceName_default = California
localityName = Locality Name (eg, city)
localityName_default = San Francisco
organizationName = Organization Name (eg, company)
organizationName_default = My Company
commonName = Common Name (e.g. server FQDN or YOUR name)
commonName_default = ${SERVER_CN}
[req_ext]
subjectAltName = @alt_names
[v3_req]
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${SERVER_CN}
EOF

# Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}" -config server_cert_ext.cnf

# Step 4: Sign the certificate with the CA we created (it's called self-signing) - server.crt
openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt -extensions v3_req -extfile server_cert_ext.cnf

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem

echo "Certificate creation completed successfully."
