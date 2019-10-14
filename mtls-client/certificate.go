package mtls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"log"
	"os"
)

func RetrieveMutualAuthCertificate(signRequest SignRequest) tls.Certificate {
	if _, err := os.Stat("client.crt"); os.IsNotExist(err) {

		certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Printf("Unable to generate RSA key: %s", err.Error())
		}

		response, err := resty.R().
			SetBody(x509.MarshalPKCS1PublicKey(&certPrivKey.PublicKey)).
			SetHeader("Content-Type", "application/json").
			SetHeader("cn", signRequest.CommonName).
			SetHeader("Authorization", signRequest.Authorization).
			Post(signRequest.CertificateAuthorityUrl + "/api/v1/sign")

		if err != nil {
			log.Fatal("Error calling CA /api/v1/sign for signing certificate, error: ", err.Error())
		}

		var signedCert ClientCertificate
		json.Unmarshal(response.Body(), &signedCert)

		//Public key
		certOut, err := os.Create("client.crt")
		if err != nil {
			log.Printf("Unable to create client cert file: %s", err.Error())
		}

		err = pem.Encode(certOut, &pem.Block{
			Type:    "CERTIFICATE",
			Bytes:   signedCert.Certificate,
			Headers: nil,
		})

		if err != nil {
			log.Printf("Unable to create client cert: %s", err.Error())
		}

		err = certOut.Close()
		if err != nil {
			log.Printf("Unable to close client file: %s", err.Error())
		}

		keyOut, err := os.OpenFile("client.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Printf("Unable to open client.key: %s", err.Error())
		}

		err = pem.Encode(keyOut, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
		})

		if err != nil {
			log.Printf("Unable to PEM encode client.key cert: %s", err.Error())
		}
		err = keyOut.Close()
		if err != nil {
			log.Printf("Unable to close client.key: %s", err.Error())
		}
	}
	return loadSignedClientCert()

}

func RetrieveCaCertificate(request SignRequest) []byte {
	if _, err := os.Stat("ca.crt"); os.IsNotExist(err) {

		resty.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })  // No CA certificate yet to verify connection
		response, err := resty.R().
			SetHeader("Content-Type", "application/json").
			Get(request.CertificateAuthorityUrl + "/api/v1/cert")

		if err != nil {
			log.Printf("Error calling CA at /api/v1/cert for CA certificate, error: %s", err.Error())
		}

		var caCert ClientCertificate

		json.Unmarshal(response.Body(), &caCert)

		//Public key
		certOut, err := os.Create("ca.crt")
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caCert.Certificate })

		certOut.Close()
	}
	return loadCACert()
}

func loadSignedClientCert() tls.Certificate {
	catls, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		panic(err)
	}
	return catls
}

func loadCACert() []byte {
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	return caCert
}

type ClientCertificate struct {
	Certificate []byte `json:"certificate"`
}

type SignRequest struct {
	CommonName              string
	CertificateAuthorityUrl string
	Authorization           string
}
