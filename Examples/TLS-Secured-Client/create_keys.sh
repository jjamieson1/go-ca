#!/bin/bash
#This is deprecated, golang does all this now.
echo "[ req ]
prompt = no
default_bits = 2048
encrypt_key = no
distinguished_name = req_distinguished_name
string_mask = utf8only

[ req_distinguished_name ]
C=CA
ST=Saskatchewan
L=Regina
O=Vivvo Application Studios
CN=localhost.vivvocloud.com
emailAddress=jamie@vivvo.com" > config.cnf

openssl genrsa -des3 -out ca.key -passout pass:Today123 4096
openssl genrsa -des3 -passout pass:Today123 -out server.key 1024
openssl req -new -passin pass:Today123 -x509 -days 365 -key ca.key -config config.cnf -out ca.crt
openssl req -new -passin pass:Today123 -config config.cnf -key server.key -out server.csr
openssl x509 -req -passin pass:Today123 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
openssl genrsa -des3 -passout pass:Today123 -out client.key 1024
openssl req -passin pass:Today123 -config config.cnf -new -key client.key -out client.csr
openssl x509 -req -passin pass:Today123 -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt
openssl rsa -passin pass:Today123 -in server.key -out temp.key
rm server.key
mv temp.key server.key
openssl rsa -passin pass:Today123 -in client.key -out temp.key
rm client.key
mv temp.key client.key
