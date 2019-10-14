package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jjamieson1/go-ca/certificate"
	"github.com/jjamieson1/go-ca/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	os.Setenv("STARTED_ON", time.Now().Format(time.RFC3339))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please create this")
	}

	cert, err := certificate.RetrieveCACertificate()
	    if err != nil {
	    	log.Printf("Error creating the CA certificate")
		}
		CA := &server.CaCertificate{CA: cert }
		certificate.CheckCreateTLSCertificate()
		tlsCert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(CA.CA.Raw)


	r := mux.NewRouter()

	r.HandleFunc("/api/v1/sign", CA.SignCertificate).Methods("POST")
	r.HandleFunc("/api/v1/cert", CA.ViewCertificate)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth: tls.NoClientCert,
		ClientCAs: caCertPool,
	}

	tlsSrv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stderr, r),
		Addr:      "0.0.0.0:443",
		TLSConfig: tlsConfig,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	//Starting HTTPS server
		log.Printf("Staring HTTPS service")
		if err := tlsSrv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
			log.Printf("Unable to start on port, error: %s", err.Error())
		}
}
