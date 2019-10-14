package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jjamieson1/go-ca/mtls-client"
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, "hello, world!\n")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/hello", HelloServer)

	signRequest := mtls.SignRequest{
		CommonName: "tls-secured-server.vivvocloud.com",
		CertificateAuthorityUrl: "https://localhost.vivvocloud.com",
		Authorization: "abc123",
	}

	caCert := mtls.RetrieveCaCertificate(signRequest)
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Printf("Unable to load CA into the cert pool")
	}

	tlsCert := mtls.RetrieveMutualAuthCertificate(signRequest)
	
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	server.ListenAndServeTLS("client.crt", "client.key") //private cert
}
