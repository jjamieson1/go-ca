package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/jjamieson1/go-ca/certificate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
	"os"
	"testing"
)

// Internal tests to test creation and generation of internal CA and TLS certs
func TestCreateCA(t *testing.T) {
	// Remove existing files
	 os.Remove("ca.crt")
	 os.Remove("ca.key")

	ca, err := certificate.RetrieveCACertificate()
	if err != nil {
	assert.AnError.Error()
	}

	if !ca.IsCA {
		err := errors.New("certificate created is not a CA")
		assert.Error(t, err)
	}
}
// External tests to test the API's
func TestGetCAFromApi(t *testing.T) {

	resty.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })  // No CA certificate yet to verify connection
	_, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(  "https://localhost/api/v1/cert")

	if err !=nil {
		assert.Error(t, err)
	}
}

func TestCreateServerTLSCertificates(t *testing.T) {
	// Remove existing  files
	os.Remove("server.crt")
	os.Remove("server.key")

	err := certificate.CheckCreateTLSCertificate()
	if err != nil {
		assert.Error(t, err)
	}
}

func TestSignClientCertificate(t *testing.T) {

	signRequest := mtls.SignRequest{
		CommonName: "tls-secured-client.vivvocloud.com",
		CertificateAuthorityUrl: "https://localhost.vivvocloud.com",
		Authorization: "abc123",
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	a := x509.MarshalPKCS1PublicKey(publicKey)
	certificateToSign := ClientCertificate{
		Certificate: a,
	}


	requestBody, err := json.Marshal(certificateToSign)

	_, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetHeader("cn", signRequest.CommonName).
		SetHeader("Authorization", signRequest.Authorization).
		Post(signRequest.CertificateAuthorityUrl + "/api/v1/sign")

	if err != nil {
		assert.Error(t, err)
	}
}

type ClientCertificate struct {
	Certificate []byte `json:"certificate"`
}

