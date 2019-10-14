# go-ca
Automation for setting up TLS mutual authentication between go services.  Setting up and maintaining TLS and mutual authentication can  be a tedious task.  This projects goal is to easily deploy a service that helps bootstrap authorized services in an installation.

This project also includes  two example projects (a client and a server), that bootstrap mutual TLS from the go-ca and perform a TLS authenticated request.  The examples also demonstrates what is required to include the client kit in your projects.
## Getting Started
By cloning this repository, you will find a file at the root of the project called .env.   This file allows you too define the environment for you installation.  I am using the good work from "github.com/joho/godotenv" to move configuration files outside of the source code.  The configuration items here define your environment.  This release is focused on a single domain PKI, signing certificates outside of a single domain is likely an enhancement at a later time.

### Prerequisites

This project is built using go modules, please use at a minimum golang 1.12. 

### Installing

Clone the repo
```
git clone https://github.com/jjamieson1/go-ca
```
cd into the directory
```
cd go-ca
```

Edit the .env file and add your parameters, then save this file in the root of the project.
```
CommonName={your fully qualified domain name}
Organization={name of your organization}
Country={two letter country code}
Province={State}
Locality={City}
StreetAddress={Address}
PostalCode={Postal Code}
```

Grab all the dependencies
```
go get -u
``` 
Run the go-ca by executing the following
```
go build . &&  go run .
```
The go-ca will start and during it's first time use it will create a self-signed root CA, and TLS certificate based on the new CA certificate.  The server should be running now on port 443.

## Testing the go-ca with the example projects.

To test locally your computer needs to be able to resolve the addresses of the three projects.  Name resolution normally refers to your local host file first, so adding your example client and server to this host file with the IP address of 127.0.0.1 will simulate this.

Here is an example of what I added to my host file.  (*note:  This is only for testing locally, not required for production*)
``` 
127.0.0.1  tls-secured-server.vivvocloud.com
127.0.0.1  tls-secured-client.vivvocloud.com
127.0.0.1  vivvo-ca.vivvocloud.com
```
Change directory to the example folders, and you will see a few example projects.

``` 
TLS-Secured-Client
TLS-Secured-Server
TLS-Secured-Server-Java
```
Change directory into the TLS-Secured-Server, and make sure there are no exisitng ca.crt, server.crt or server.key files.  Delete them if they exist.

Check out the main.go file, and you will see:

``` 
signRequest := mtls.SignRequest{
		CommonName: "tls-secured-server.vivvocloud.com",
		CertificateAuthorityUrl: "https://vivvo-ca.vivvocloud.com",
		Authorization: "abc123",
	}
```
Edit this file and change it to the domain name of the go-ca.

Save this file and start the server with:

```` 
go build . && go run .
````

The  server will:
    * Download the CA from the go-ca server and install it in it's trust store.
    * Create a private key and send a signing request to the go-ca.
    * The go ca will wign the request and send back a signed cert.
    * The server will install this cert into it's trust store.
    
The server will now be running on it's test port and will only serve requests to clients that have a signed certificate from the g-ca. 

The "tls.RequireAndVerifyClientCert" extra option on the tls-config makes this possible:

``` 
tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
```
Next change directory to the TLS-Secured-Client and take a look at the main.go,  do the same thing you did to the server client request object and modify it to match your domainname.

```
	signRequest := mtls.SignRequest{
		CommonName: "tls-secured-client.vivvocloud.com",
		CertificateAuthorityUrl: "https://localhost.vivvocloud.com",
		Authorization: "abc123",
	}
```

Save this file and start the client with:

```` 
go build . && go run .
````

A successful test should display that famous "hello world" string.

## Deployment

Deployment into a live environment is a matter of looking at the individual go applications that participate in your enterprise, and adding a dependancy in imports.

```  
import (
	"github.com/jjamieson1/go-ca/mtls-client"
)   
```

And adding the following lines to your project, and edit the "CommonName and CertificateAuthorityUrl" to bootstrap the certificates required for mutual TLS authentication

Example:
``` 
signRequest := mtls.SignRequest{
		CommonName: "tls-secured-client.vivvocloud.com",
		CertificateAuthorityUrl: "https://localhost.vivvocloud.com",
		Authorization: "abc123",
	}

	caCert := mtls.RetrieveCaCertificate(signRequest)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsCert := mtls.RetrieveMutualAuthCertificate(signRequest)

```

Configure the TLS settings:

Client example:

``` 
tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      caCertPool,
		ClientCAs: caCertPool,
	}
tlsConfig.BuildNameToCertificate()
```

Server Example:

``` 
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
tlsConfig.BuildNameToCertificate()
```

For a complete example, please refer to the main.go in either example project.

## Built With

* [golang cryto](https://golang.org/pkg/crypto/) - Golang crypto libraries


## Contributing

Please feel free to add features or fixes and issue a pull request.

## Versioning

1.0
## Authors

* **Jamie Jamieson** - *Initial work* - [jjamieson1](https://github.com/jjamieson1)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments