# Example Client using the mtls go-ca security
## How it functions
* When the service first runs it will attempt to load the CA into the cert pool.  If the file ca.crt does not exist, we call out to the go-ca service and grab the CA.
* The next thing the service does, it looks to see if there are any local client certificates created.  If there is none, a rsa 4096 keypair is created, and the public key is sent to the go-ca for signing.  The go-ca will return a signed certificate that is required for tls mutual authentication. 
* In the main method check ou the example on how to configure tls to use these certificates.
## Troubleshooting
1.  Turn it off and rm *.crt *.key files, and turn it back on.  This is only necessary if the CA has changed.
2.  Useful commands to verify certificates are:
* Verify the cert was signed by the CA:
```
openssl verify -verbose -CAfile ca.crt client.crt
```
* Verify the certificate 

```
openssl x509 -in client.crt -text -noout  
```

* Use the cert and CA against a server that uses the CA's mutual auth


``` 
 cat client.key > client.pem
 cat client.crt >> client.pem
 curl -v --cacert ca.crt --cert client.pem https://localhost.vivvocloud.com:8080/hello

```
# Generate certificate manually (deprecated)
First of all I need to generate SSL certificates to client and server. I’m creating my own certificate authority(CA) to issue the certificates.

Generate CA Certificate and Key
Generating CA Certificate/Key use to issue/sign the server and client certificates.

## CA key and certificate
```
openssl genrsa -des3 -out ca.key 4096
openssl req -new -x509 -days 365 -extensions v3_ca -key ca.key -out ca.crt
```
Generate server Key and CSR
Need to generate server key and certificate signing request(CSR) to obtain the server certificate.

## server key
```
openssl genrsa -des3 -out server.key 1024
```
## CSR (certificate sign request) to obtain certificate
```
openssl req -new -key server.key -out server.csr
```
## Generate server certificate
Generated certificate signing request(CSR) need to be signed by CA’s certificate/key to obtain the server certificate. It’s a self signed certificate.

## sign server CSR with CA certificate and key
```
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
```
## Generate client Key and CSR
Need to generate client key and certificate signing request(CSR) to obtain the client certificate.

## client key
```
openssl genrsa -des3 -out client.key 1024
```
## CSR to obtain certificate
```
openssl req -new -key client.key -out client.csr
```
## Generate client certificate
Generated certificate signing request(CSR) need to be signed by CS’s certificate/key to obtain the client certificate. It’s a self sign certificate.

## sign client CSR with CA certificate and key
```
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt
```
##Remove pass phrase from server and client keys
When generating server key and client key, its asking for a password. We need to specify that password when loading the key(in nginx and golang http client). Otherwise it will give an error. I’m removing that pass phrase from the key, then I can use the key without the password.

In here I’m copying the content of server and client key to temp file and making the temp as the key

## server key out to temp.key
```
openssl rsa -in server.key -out temp.key
```
## remove server.key
```
rm server.key
```
## make temp.key as server key
```
mv temp.key server.key
```
By same way I’m removing the pass phrase of the client key.

## client key out to temp.key
```
openssl rsa -in client.key -out temp.key
```
## remove client.key
```
rm client.key
```
### make temp.key as client key
```
mv temp.key client.key
```
