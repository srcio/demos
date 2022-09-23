package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
)

// You'll need a user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func main() {

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: "me@srcio.cn",
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	// config.CADirURL = "http://192.168.99.100:4000/directory"
	config.CADirURL = lego.LEDirectoryProduction
	config.Certificate.KeyType = certcrypto.RSA4096

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	// err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	// client.Challenge.SetHTTP01Provider(http01.NewProviderServer(lego.LEDirectoryProduction, "80"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer(lego.LEDirectoryProduction, "443"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// You should set the envs for your provider.
	{
		os.Setenv("ALICLOUD_ACCESS_KEY", "xxxxxx")
		os.Setenv("ALICLOUD_SECRET_KEY", "xxxxxx")
	}

	p, _ := dns.NewDNSChallengeProviderByName("alidns")
	client.Challenge.SetDNS01Provider(p)

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"mydomain.com", "www.mydomain.com"},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("certificates.CertURL: %v\n", certificates.CertURL)
	fmt.Printf("certificates.CertStableURL: %v\n", certificates.CertStableURL)
	fmt.Printf("certificates.CSR: %v\n", string(certificates.CSR))
	fmt.Printf("certificates.Domain: %v\n", certificates.Domain)
	fmt.Printf("certificates.Certificate: %v\n", string(certificates.Certificate))
	fmt.Printf("certificates.IssuerCertificate: %v\n", string(certificates.IssuerCertificate))
	fmt.Printf("certificates.PrivateKey: %v\n", string(certificates.PrivateKey))

	// ... all done.
}
