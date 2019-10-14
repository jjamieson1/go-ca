package e2e

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Vivvo/go-sdk/mtls"
	"github.com/Vivvo/vivvo-ca/certificate"
	"github.com/Vivvo/vivvo-ca/server"
	"gopkg.in/resty.v1"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestE2E(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	os.Remove("ca.crt")
	os.Remove("ca.key")
	os.Remove("server.crt")
	os.Remove("server.key")
	os.Remove("client.crt")
	os.Remove("client.key")

	mockCAServer := bootstrapMockCAServer()
	defer mockCAServer.Close()

	mockServer := bootstrapTLSSecuredServer(mockCAServer)
	defer mockServer.Close()

	signRequest := mtls.SignRequest{
		CommonName:              mockServer.URL,
		CertificateAuthorityUrl: mockCAServer.URL,
		Authorization:           "abc123",
	}

	caCert := mtls.RetrieveCaCertificate(signRequest)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:   caCertPool,
		ClientCAs: caCertPool,
		InsecureSkipVerify: true,
	}
	resty.SetTLSClientConfig(tlsConfig)

	tlsCert := mtls.RetrieveMutualAuthCertificate(signRequest)
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      caCertPool,
		//ClientCAs:    caCertPool,
		//InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	resty.SetTLSClientConfig(tlsConfig)

	response, err := resty.R().
		Get(mockServer.URL)

	if err != nil {
		log.Printf("Error calling CA for CA certificate, error: %s", err.Error())
	}

	log.Printf("Returned: %s", response.Body())
}

func bootstrapTLSSecuredServer(mockCAServer *httptest.Server) *httptest.Server {
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("NICE WORK! YOU'RE TLS'd UP"))
	}))
	signRequest := mtls.SignRequest{
		CommonName:              mockServer.URL,
		CertificateAuthorityUrl: mockCAServer.URL,
		Authorization:           "abc123",
	}
	caCert := mtls.RetrieveCaCertificate(signRequest)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:   caCertPool,
		ClientCAs: caCertPool,
	}
	resty.SetTLSClientConfig(tlsConfig)

	log.Println("Made it!")
	tlsCert := mtls.RetrieveMutualAuthCertificate(signRequest)

	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	mockServer.TLS = tlsConfig
	mockServer.StartTLS()
	return mockServer
}

func bootstrapMockCAServer() ( *httptest.Server) {
	ca, _ := certificate.RetrieveCACertificate()
	VivvoCA := &server.CaCertificate{CA: ca}
	certificate.CheckCreateTLSCertificate()
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(VivvoCA.CA.Raw)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sign", VivvoCA.SignCertificate).Methods("POST")
	r.HandleFunc("/api/v1/cert", VivvoCA.ViewCertificate)


	c, _ := tls.LoadX509KeyPair("./server.crt", "./server.key")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{c},
		ClientAuth: tls.NoClientCert,
		ClientCAs:  caCertPool,
	}
	mockCAServer := httptest.NewUnstartedServer(r)
	mockCAServer.TLS = tlsConfig
	mockCAServer.StartTLS()
	return mockCAServer
}
