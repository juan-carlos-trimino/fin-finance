package security

/***
For more information about the openssl req command options, visit the OpenSSL req documentation
page (https://www.openssl.org/docs/man1.0.2/man1/openssl-req.html).

To check the version of OpenSSL, run the following command:
$ openssl version -a

The first step to obtain an SSL certificate is using OpenSSL to create a Certificate Signing
Request (CSR) that can be sent to a Certificate Authority (CA). CSR is a block of encoded text
with data about the website and the company. In order for a CSR to be created, it needs to have a
private key from which the public key is extracted. While an existing key can be used, it's
recommended to always generate a new private key whenever a CSR is created.
$ openssl genrsa -out priv_server.key 2048

or

if using a passphrase
$ openssl genrsa -out priv_server.key 2048 -passout pass:xyzxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Please note that certain servers will not accept private keys with passphrases.

To view the encoded contents of the private key, run the following command:
$ cat priv_server.key

To check the private key, run the command below:
$ openssl rsa -in priv_server.key -check

Use the following command to decode the private key and view its contents:
$ openssl rsa -text -in priv_server.key -noout
The -noout switch omits the output of the encoded version of the private key.

The private key file contains both the private key and the public key. To extract the public key
from the private key file, run the command below:
$ openssl rsa -in priv_server.key -pubout -out pub_server.key

After generating the private key, the CSR can be created. The CSR is created using the PEM format,
and it contains the public key portion of the private key as well as information about the company.
Run the following command to generate the CSR:
$ openssl req -new -key priv_server.key -out priv_server.csr \
-subj "/C=US/ST=TX/L=Houston/O=Company Name/OU=Dept A/CN=localhost"

Instead of generating a private key and then creating a CSR in two separate steps, both steps can
actually be performed one step. Use the following command to create both the private key and CSR:
$ openssl req -newkey rsa:2048 -nodes -keyout priv_server.key -out priv_server.csr \
-subj "/C=US/ST=TX/L=Houston/O=Company Name/OU=Dept A/CN=localhost"

After creating the CSR using the private key, it is recommended to verify that the information
contained in the CSR is correct and that the file hasn't been modified or corrupted. Use the
following command to view the information in the CSR before submitting it to the CA:
$ openssl req -text -noout -verify -in priv_server.csr
The -noout switch omits the output of the encoded version of the CSR. The -verify switch checks the
signature of the file to make sure it hasn't been modified.

To send the CSR to the CA, view the raw output of the CSR by running the command below:
$ cat priv_server.csr
Make sure to include the -----BEGIN CERTIFICATE REQUEST----- and -----END CERTIFICATE REQUEST-----
tags and paste everything into the SSL vendor's order form.

After receiving the certificate from the CA, it is recommended to ensure that the information in
the certificate is correct and matches the private key. Use the following command to view the
contents of the certificate:
$ openssl x509 -text -in priv_server.crt -noout

---------------------------------------------------------------------------------------------------
In cryptography and computer security, self-signed certificates are public key certificates that
are not issued by a CA. These self-signed certificates are easy to create and are free. However,
they do not provide any trust value.
---------------------------------------------------------------------------------------------------
To generate a self-signed SSL certificate using OpenSSL, complete the following steps:
$ openssl req -x509 -newkey rsa:2048 -keyout rootCA.key -out rootCA.crt \
-sha512 -days 1 -nodes \
-subj "/C=US/ST=TX/L=Houston/O=Company Name/OU=Dept A/CN=localhost"

Once the certificate has been generated, verify that it is correct as per the parameters that were
used.
$ openssl x509 -text -noout -in rootCA.crt
***/

import (
	// "bytes"
	// "bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"net"

	// "crypto/rsa"
	// "crypto/sha256"
	"crypto/rsa"
	// "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
	// "crypto/x509"
	// "crypto/x509/pkix"
	// "encoding/base64"
	// "encoding/pem"
	// "math/big"
	// "time"
	// "github.com/pkg/errors"
)

func genCert(template, parent *x509.Certificate, pubKey *rsa.PublicKey,
             parentPrivKey *rsa.PrivateKey) (*x509.Certificate, []byte) {
  /***
  The certificate is signed by parent. If parent is equal to template then the certificate is
  self-signed. The parameter pubKey is the public key of the certificate to be generated and
  parentPrivKey is the private key of the signer.
  ***/
  certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pubKey, parentPrivKey)
  if err != nil {
    panic("Failed to create certificate.\n" + err.Error())
  }
  cert, err := x509.ParseCertificate(certBytes)
  if err != nil {
    panic("Failed to parse certificate.\n" + err.Error())
  }
  //PEM encode (https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail) the certificate.
  certPEM := pem.EncodeToMemory(&pem.Block{
    Type: "CERTIFICATE",
    Bytes: certBytes,
  })
  return cert, certPEM
}

func rootCATemplate() x509.Certificate {
  //Create the Certificate Request (CR). The CR will be signed with the private key; this provides
  //the identity of the requester.
  template := x509.Certificate{
    /***
    The serial number is an integer assigned by the CA to each certificate. It MUST be unique for
    each certificate issued by a given CA; i.e., the issuer name and serial number identify a
    unique certificate.
    ***/
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"RootCA, Inc."},
      Country: []string{"US"},
      Province: []string{"Texas"},
      Locality: []string{"Houston"},
      StreetAddress: []string{"721 Tree St."},
      PostalCode: []string{"77909"},
      //CommonName: "Root CA",
    },
    NotBefore: time.Now(),
    NotAfter: time.Now().AddDate(0, 0, 7),
    IsCA: true,  //A CA certificate.
    BasicConstraintsValid: true,
    //Specify how the certificate's public key may be used.
    //Subject public key is used to verify signatures on certificates.
    //This extension must only be used for CA certificates.
    KeyUsage: x509.KeyUsageCertSign |
    //Certificate may be used to apply a digital signature.
              x509.KeyUsageDigitalSignature |
    //Subject public key is to verify signatures on revocation information, such as a CRL.
    //This extension must only be used for CA certificates
              x509.KeyUsageCRLSign,
    //CAs/ICAs should not have any EKUs specified
    ExtKeyUsage: []x509.ExtKeyUsage{},
//    MaxPathLen: 2,
//    IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
//    DNSNames: []string{"localhost"},
  }
  return template
}

func GenRootCA() (*x509.Certificate, []byte, *rsa.PrivateKey) {
  template := rootCATemplate()
  /***
  A Key Pair consists of a Private and Public key. The private key must be secured; the public key
  is derived from the private key and can be shared. Before generating a certificate, a Key Pair is
  needed.
  ***/
  privKey, err := rsa.GenerateKey(rand.Reader, 2048)  //Generate a RSA keypair.
  if err != nil {
    panic("Failed to generate a RSA keypair.\n" + err.Error())
  }
  //Create a self-signed certificate.
  cert, certPEM := genCert(&template, &template, &privKey.PublicKey, privKey)
  return cert, certPEM, privKey
}

func intermediateCATemplate() x509.Certificate {
  template := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"Intermediate CA, Inc."},
      Country: []string{"US"},
      Province: []string{"Texas"},
      Locality: []string{"Houston"},
      StreetAddress: []string{"721 Tree St."},
      PostalCode: []string{"77909"},
      //CommonName: "SelfTLS CA",
    },
    NotBefore: time.Now(),
    NotAfter: time.Now().AddDate(0, 0, 7),
    IsCA: true,
    BasicConstraintsValid: true,
    KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageCRLSign,
    //CAs/ICAs should not have any EKUs specified
    ExtKeyUsage: []x509.ExtKeyUsage{},
    // MaxPathLenZero: false,
    // MaxPathLen: 1,
    // IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
  }
  return template
}

func GenIntermediateCA(cert *x509.Certificate, certPrivKey *rsa.PrivateKey) (*x509.Certificate,
                       []byte, *rsa.PrivateKey) {
  template := intermediateCATemplate()
  interPrivKey, err := rsa.GenerateKey(rand.Reader, 2048);
  if err != nil {
    panic(err)
  }
  interCert, interCertPEM := genCert(&template, cert, &interPrivKey.PublicKey, certPrivKey)
  return interCert, interCertPEM, interPrivKey
}

func serverTemplate(hosts []string) x509.Certificate {
  template := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"Company, Inc."},
      Country: []string{"US"},
      Province: []string{"Texas"},
      Locality: []string{"Dallas"},
      StreetAddress: []string{"3000 Lake Dr"},
      PostalCode: []string{"76092"},
      CommonName: "localhost",
    },
    NotBefore: time.Now(),
    NotAfter: time.Now().AddDate(0, 0, 7),  //Valid for seven day.
    IsCA: false,
    BasicConstraintsValid: true,
    KeyUsage: x509.KeyUsageDigitalSignature |
    //Certificate enables use of a key agreement protocol to establish a symmetric key with a target.
    //Symmetric key may then be used to encrypt & decrypt data sent between the entities.
              x509.KeyUsageKeyAgreement |
    //Certificate may be used to encrypt a symmetric key which is then transferred to the target.
    //Target decrypts key, subsequently using it to encrypt & decrypt data between the entities.
              x509.KeyUsageKeyEncipherment,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    // MaxPathLenZero: true,
    // IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
    // DNSNames: []string{"localhost"},
  }
  //hosts: []string{"hostname", "ipv4addr", "ipv6addr"}
  for _, h := range hosts {
    if ip := net.ParseIP(h); ip != nil {
      template.IPAddresses = append(template.IPAddresses, ip)
    } else {
      template.DNSNames = append(template.DNSNames, h)
      template.Subject.CommonName = h
    }
  }
  return template
}

func GenServerCert(caCert *x509.Certificate, caCertPrivKey *rsa.PrivateKey) (*x509.Certificate,
                   []byte, []byte) {
  template := serverTemplate([]string{"localhost", "::1", "127.0.0.1"})
  serverPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    panic("Failed to generate the server key.\n" + err.Error())
  }
  serverCert, serverCertPEM := genCert(&template, caCert, &serverPrivKey.PublicKey, caCertPrivKey)
  serverPrivKeyPEM := pem.EncodeToMemory(&pem.Block{
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey),
  })
  return serverCert, serverCertPEM, serverPrivKeyPEM
}

func VerifyIntermediateCA(rootCert, interCert *x509.Certificate) {
  roots := x509.NewCertPool()
  roots.AddCert(rootCert)
  opts := x509.VerifyOptions{
    Roots: roots,
  }
  //
  if _, err := interCert.Verify(opts); err != nil {
    panic("Failed to verify certificate.\n" + err.Error())
  }
  fmt.Println("Intermediate CA verified.")
}

func VerifyCertificateChain(rootCert, interCert, serverCert *x509.Certificate) {
  roots := x509.NewCertPool()
  roots.AddCert(rootCert)
  inter := x509.NewCertPool()
  inter.AddCert(interCert)
  opts := x509.VerifyOptions{
    Roots: roots,
    Intermediates: inter,
    // Intermediates: x509.NewCertPool(),
  }
  //
  if _, err := serverCert.Verify(opts); err != nil {
    panic("Failed to verify certificate chain.\n" + err.Error())
  }
  fmt.Println("Certificate chain is valid.")
}

func KeyToPemBlock(key interface{}) *pem.Block {
  switch k := key.(type) {
  case *rsa.PrivateKey:
    return &pem.Block{
      Type: "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(k),
    }
  case *ecdsa.PrivateKey:
    if b, err := x509.MarshalECPrivateKey(k); err != nil {
      panic("Unable to marshal ECDSA private key.\n" + err.Error())
    } else {
      return &pem.Block{
        Type: "EC PRIVATE KEY",
        Bytes: b,
      }
    }
  default:
    return nil
  }
}

/***

// GenerateSelfSignedCert creates a self-signed certificate and key for the given host.
// Host may be an IP or a DNS name
// The certificate will be created with file mode 0644. The key will be created with file mode 0600.
// If the certificate or key files already exist, they will be overwritten.
// Any parent directories of the certPath or keyPath will be created as needed with file mode 0755.
func GenerateSelfSignedCert(host, certPath, keyPath string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, time.Now().Unix()),
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}

	// Write cert
	if err := os.MkdirAll(filepath.Dir(certPath), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(certPath, certBuffer.Bytes(), os.FileMode(0644)); err != nil {
		return err
	}

	// Write key
	if err := os.MkdirAll(filepath.Dir(keyPath), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(keyPath, keyBuffer.Bytes(), os.FileMode(0600)); err != nil {
		return err
	}

	return nil
}
----------------------
package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// GenerateSelfSignedCert creates a self-signed certificate and key for the given host.
// Host may be an IP or a DNS name
// The certificate will be created with file mode 0644. The key will be created with file mode 0600.
// If the certificate or key files already exist, they will be overwritten.
// Any parent directories of the certPath or keyPath will be created as needed with file mode 0755.
func GenerateSelfSignedCert(host, certPath, keyPath string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, time.Now().Unix()),
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
  var b bytes.Buffer // A Buffer needs no initialization.
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}

	// Write cert
	if err := os.MkdirAll(filepath.Dir(certPath), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(certPath, certBuffer.Bytes(), os.FileMode(0644)); err != nil {
		return err
	}

	// Write key
	if err := os.MkdirAll(filepath.Dir(keyPath), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(keyPath, keyBuffer.Bytes(), os.FileMode(0600)); err != nil {
		return err
	}

	return nil
}
***/

