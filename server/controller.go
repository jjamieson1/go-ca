package server

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/rsa"
	"crypto/x509"
	"github.com/jjamieson1/go-ca/certificate"
	"io/ioutil"
	"log"
	"net/http"
)

func (ca *CaCertificate) ViewCertificate(w http.ResponseWriter, r *http.Request) {
	c := certificate.ClientCertificate{Certificate: ca.CA.Raw	}
	WriteJSON(c, http.StatusOK, w)
}

//Todo add ACI to this endpoint


func (ca *CaCertificate) SignCertificate(w http.ResponseWriter, r *http.Request) {
	//Todo add ACI to this endpoint


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("No public key supplied, proceeding to create key: %s", err)
		return
	}

	if len(body) == 0 {
		certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Printf("Unable to generate RSA key: %s", err.Error())
		}
		body = x509.MarshalPKCS1PublicKey(&certPrivKey.PublicKey)
	}

	WriteJSON(certificate.SignClientCertificateRequest(r.Header.Get("cn"), body), http.StatusOK, w)
}

type CaCertificate struct {
	CA *x509.Certificate
}
