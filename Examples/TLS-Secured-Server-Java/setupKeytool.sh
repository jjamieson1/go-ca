#!/bin/bash
export CA=localhost
export CN=tls-java-test.vivvocloud.com
export PORT=443
keytool -genkey -noprompt -keypass password -storepass password -keyalg RSA -alias client -dname "CN=Vivvo, OU=ID1, O=Vivvo, L=Regina, ST=Saskatchewan, C=CA" -keystore keystore.jks -storepass password -validity 3600 -keysize 2048
curl --insecure https://${CA}:${PORT}/api/v1/cert -o ca.json
cat ca.json | jq -r '.certificate' | base64 --decode > ca.crt
keytool -noprompt -storepass password -import -trustcacerts -alias root -file ca.crt -keystore keystore.jks
curl --insecure -X POST --header "cn: ${CN}" https://${CA}:${PORT}/api/v1/sign -o signed.json
cat signed.json | jq -r '.certificate' | base64 --decode > certificate.crt
keytool -noprompt -storepass password -import -trustcacerts -alias server -file certificate.crt -keystore keystore.jks
keytool -noprompt -storepass password -list -keystore keystore.jks
rm ca.json signed.json
