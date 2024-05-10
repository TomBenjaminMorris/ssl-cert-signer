package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// Helpers
func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

func checkIfValidFile(filename string) (bool, error) {
	// Check if file is CSR
	if fileExtension := filepath.Ext(filename); fileExtension != ".csr" {
		return false, fmt.Errorf("file %s is not a CSR", filename)
	}

	// Check if file does exist
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", filename)
	}

	return true, nil
}

//go:embed ca_files/rootCA*
var ca_files embed.FS

type CA struct {
	ca  *x509.Certificate
	key *rsa.PrivateKey
}

func loadCA() (CA, error) {
	// Load the CA certificate
	certPEMBytes, err := ca_files.ReadFile("ca_files/rootCA.crt")
	check(err)

	block, _ := pem.Decode(certPEMBytes)
	if block == nil {
		log.Fatalf("err %v", err)
	}

	ca, err := x509.ParseCertificate(block.Bytes)
	check(err)

	keyPEMBytes, err := ca_files.ReadFile("ca_files/rootCA.key")
	check(err)

	privPem, _ := pem.Decode(keyPEMBytes)
	parsedKey, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
	privkey := parsedKey.(*rsa.PrivateKey)
	check(err)

	certificateAuthority := CA{ca: ca, key: privkey}

	return certificateAuthority, nil
}

type inputFile struct {
	filepath string
	flag1    bool
	flag2    bool
}

func getCSR() (inputFile, error) {
	// Validate arguments
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("filepath argument is required")
	}

	flag1 := flag.Bool("flag1", false, "Placeholder flag1")
	flag2 := flag.Bool("flag2", false, "Placeholder flag2")

	flag.Parse()
	fileLocation := flag.Arg(0)

	return inputFile{fileLocation, *flag1, *flag2}, nil
}

type CSR struct {
	csr *x509.CertificateRequest
}

func parseCSR(input inputFile) CSR {
	// Load the CSR
	csrBytes, err := os.ReadFile(input.filepath)
	check(err)

	block, _ := pem.Decode(csrBytes)
	if block == nil {
		log.Fatalf("err %v", err)
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	check(err)

	signingRequest := CSR{csr}

	return signingRequest
}

func signCSR(input CSR, ca CA) []byte {
	// sign the CSR
	csr_notBefore := time.Now()
	csr_notAfter := csr_notBefore.Add(time.Hour * 24 * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	csr_serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	check(err)

	ccr := &x509.Certificate{
		SerialNumber: csr_serialNumber,
		Subject: pkix.Name{
			Organization:       input.csr.Subject.Organization,
			OrganizationalUnit: input.csr.Subject.OrganizationalUnit,
			Locality:           input.csr.Subject.Locality,
			Province:           input.csr.Subject.Province,
			Country:            input.csr.Subject.Country,
			CommonName:         input.csr.Subject.CommonName,
		},
		NotBefore:             csr_notBefore,
		NotAfter:              csr_notAfter,
		DNSNames:              input.csr.DNSNames,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	cert, err := x509.CreateCertificate(rand.Reader, ccr, ca.ca, input.csr.PublicKey, ca.key)
	check(err)

	return cert
}

func writeCert(cert []byte) {
	filepath := "exampleFiles/mydomain.com.crt"
	certOut, err := os.Create(filepath)
	check(err)

	if err := pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
		log.Fatalf("Failed to write data: %s", err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
		log.Fatalf("Failed to write data: %s", err)
	}

	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing %s  %s", filepath, err)
	}
}
