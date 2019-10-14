package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func RetrieveCACertificate() (*x509.Certificate, error) {
	caCertificate := &x509.Certificate{}
	var err error
	if _, err := os.Stat("ca.crt"); os.IsNotExist(err) {
		log.Printf("Creating a new CA for this environment")
		ca := &x509.Certificate{
			SerialNumber: big.NewInt(1653),
			Subject: pkix.Name{
				CommonName: os.Getenv("CommonName"),
				Organization:  []string{os.Getenv("Organization")},
				Country:       []string{os.Getenv("Country")},
				Province:      []string{os.Getenv("Province")},
				Locality:      []string{os.Getenv("Locality")},
				StreetAddress: []string{os.Getenv("StreetAddress")},
				PostalCode:    []string{os.Getenv("PostalCode")},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		priv, _ := rsa.GenerateKey(rand.Reader, 2048)
		pub := &priv.PublicKey
		ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
		if err != nil {
			log.Println("create ca failed", err.Error())
		}

		// Public key
		certOut, err := os.Create("ca.crt")
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
		certOut.Close()

		// Private key
		keyOut, err := os.OpenFile("ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
		keyOut.Close()

	}

	caCertificate = loadCA()

	return caCertificate, err
}

func CheckCreateTLSCertificate() error {
	var err error
	if _, err := os.Stat("server.crt"); os.IsNotExist(err) {
		catls, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
		if err != nil {
			return err
		}
		ca, err := x509.ParseCertificate(catls.Certificate[0])
		if err != nil {
			return err
		}
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(1658),

			Subject: pkix.Name{
				CommonName: 	os.Getenv("CommonName"),
				Organization:  []string{os.Getenv("Organization")},
				Country:       []string{os.Getenv("Country")},
				Province:      []string{os.Getenv("Province")},
				Locality:      []string{os.Getenv("Locality")},
				StreetAddress: []string{os.Getenv("StreetAddress")},
				PostalCode:    []string{os.Getenv("PostalCode")},
			},
			DNSNames: []string{"*.vivvocloud.com"},
			NotBefore:    time.Now(),
			NotAfter:     time.Now().AddDate(10, 0, 0),
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			//ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		}
		priv, _ := rsa.GenerateKey(rand.Reader, 2048)
		pub := &priv.PublicKey
		cert_b, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)
		certOut, err := os.Create("server.crt")
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert_b})
		certOut.Close()
		log.Print("written cert.pem\n")

		// Private key
		keyOut, err := os.OpenFile("server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
		keyOut.Close()
		log.Print("written key.pem\n")

	}
	return err
}

func SignClientCertificateRequest(commonName string, publicKey []byte) ClientCertificate {

	catls, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		panic(err)
	}

	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		log.Printf("failed to parse CA certificate: " + err.Error())
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:    commonName,
			Organization:  []string{os.Getenv("Organization")},
			Country:       []string{os.Getenv("Country")},
			Province:      []string{os.Getenv("Province")},
			Locality:      []string{os.Getenv("Locality")},
			StreetAddress: []string{os.Getenv("StreetAddress")},
			PostalCode:    []string{os.Getenv("PostalCode")},
		},
		//IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")}, //Enable this for the E2E tests
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	b, err := x509.ParsePKCS1PublicKey(publicKey)

    signedCert, err := x509.CreateCertificate(rand.Reader, cert, ca, b, catls.PrivateKey)
    if err != nil {
    	log.Printf("Error signing the certificate, Error: %s", err.Error())
	}
	clientCertificate := ClientCertificate{
				Certificate: signedCert,
			}
	log.Printf("Created and returned a certificate for %s", commonName)
	return clientCertificate
}

func loadCA() *x509.Certificate{
	catls, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}
	return ca
}

type ClientCertificate struct {
	Certificate 	[]byte	`json:"certificate"`
}