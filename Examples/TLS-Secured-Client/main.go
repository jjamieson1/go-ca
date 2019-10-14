package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jjamieson1/go-ca/mtls-client"
	"gopkg.in/resty.v1"
	"log"
)

func main() {

	signRequest := mtls.SignRequest{
		CommonName: "tls-secured-client.vivvocloud.com",
		CertificateAuthorityUrl: "https://localhost.vivvocloud.com",
		Authorization: "abc123",
	}

	caCert := mtls.RetrieveCaCertificate(signRequest)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsCert := mtls.RetrieveMutualAuthCertificate(signRequest)


	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      caCertPool,
		ClientCAs: caCertPool,

	}
	tlsConfig.BuildNameToCertificate()

		resty.SetTLSClientConfig(tlsConfig)
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get("https://tls-secured-server.vivvocloud.com:8080/hello")

	if err != nil {
		log.Printf("Error calling CA for CA certificate, error: %s", err.Error())
	}

	log.Printf("Returned: %s", response.Body())

}
