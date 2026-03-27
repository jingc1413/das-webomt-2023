package certgen

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	_ "embed"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/mail"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

//go:embed cert
var certFS embed.FS

var (
	expiredDaysCA = 365 * 30 // days
	rootCaName    = "rootCA.pem"
	rootKeyName   = "rootCA-key.pem"

	keyFile, certFile, p12File string

	userAndHostname string

	certGenOnce      sync.Once
	certGenSingleton *CertGen
)

type CertGenOptions struct {
	Hosts []string //use by gen cert. Subject Alternate Name values. default host interfaces

	Organization     string //default empty
	OrganizationUnit string //default empty
	CommonName       string //default empty

	Country                   string //default empty
	Locality, Province        string //default empty
	StreetAddress, PostalCode string //default empty

	ExpirationTime int // unit: day, default 365*20
}

type CertResult struct {
	CaPEM   []byte
	CertPEM []byte
	PrivPEM []byte
}

func GetCertResult(option *CertGenOptions) (*CertResult, error) {
	s := MustGetCertGen()
	{
		o := option
		if o == nil {
			o = &CertGenOptions{}
		}
		if o.ExpirationTime < 1 {
			o.ExpirationTime = 20 * 365
		}
		if len(o.CommonName) > 0 {
			o.Hosts = append(o.Hosts, o.CommonName)
		}
		s.certOption = option
	}
	if e := s.makeCertToMem(s.certOption.Hosts); e != nil {
		return nil, errors.Wrapf(e, "make cert to mem")
	}
	return &CertResult{
		CaPEM:   s.caPEM,
		CertPEM: s.certPEM,
		PrivPEM: s.privPEM,
	}, nil
}

func MustGetCertGen() *CertGen {
	certGenOnce.Do(func() {
		u, err := user.Current()
		if err == nil {
			userAndHostname = u.Username + "@"
		}
		if h, err := os.Hostname(); err == nil {
			userAndHostname += h
		}
		if err == nil && u.Name != "" && u.Name != u.Username {
			userAndHostname += " (" + u.Name + ")"
		}
		if certGenSingleton, err = NewCertGen(); err != nil {
			panic("init gen cert")
		}
	})
	if certGenSingleton == nil {
		panic("setup cert failed")
	}
	return certGenSingleton
}

//func GetCertBlock() ([]byte, []byte) {
//	s := MustGetCertGen()
//	return s.GetCertBlock()
//}

//func GetCaBlock() []byte {
//	s := MustGetCertGen()
//	return s.caPEM
//}

type CertGen struct {
	pkcs12, ecdsa, client bool

	CAROOT string

	caCert *x509.Certificate
	caKey  crypto.PrivateKey

	caPEM    []byte
	caKeyPEM []byte
	certPEM  []byte
	privPEM  []byte

	certOption *CertGenOptions
}

func NewCertGen() (*CertGen, error) {
	s := &CertGen{}
	s.CAROOT = getCAROOT()
	e := s.loadFromFS()
	if e != nil {
		return nil, errors.Wrapf(e, "load ca from fs")
	}
	return s, nil
}

func (m *CertGen) makeCertToMem(hosts []string) error {
	if hosts == nil {
		ips, e := getAllIps()
		if e != nil {
			return errors.Wrapf(e, "get all ips")
		}
		hosts = ips
	}
	if m.caKey == nil {
		return errors.New("can't create certificates, ca-key is missing")
	}
	priv, err := m.generateKey(false)
	if err != nil {
		return errors.Wrapf(err, "generate certificate key")
	}
	pub := priv.(crypto.Signer).Public()

	expiration := time.Now().AddDate(0, 0, m.certOption.ExpirationTime)
	sn, err := m.randomSerialNumber()
	if err != nil {
		return err
	}
	tpl := &x509.Certificate{
		Version:      3,
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization:       []string{m.certOption.Organization},
			OrganizationalUnit: []string{m.certOption.OrganizationUnit},
			CommonName:         m.certOption.CommonName,
			Country:            []string{m.certOption.Country},
			Province:           []string{m.certOption.Province},
			Locality:           []string{m.certOption.Locality},
			StreetAddress:      []string{m.certOption.StreetAddress},
			PostalCode:         []string{m.certOption.PostalCode},
		},
		NotBefore: time.Now(),
		NotAfter:  expiration,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			tpl.IPAddresses = append(tpl.IPAddresses, ip)
		} else if email, err := mail.ParseAddress(h); err == nil && email.Address == h {
			tpl.EmailAddresses = append(tpl.EmailAddresses, h)
		} else if uriName, err := url.Parse(h); err == nil && uriName.Scheme != "" && uriName.Host != "" {
			tpl.URIs = append(tpl.URIs, uriName)
		} else {
			tpl.DNSNames = append(tpl.DNSNames, h)
		}
	}
	if m.client {
		tpl.ExtKeyUsage = append(tpl.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}
	if len(tpl.IPAddresses) > 0 || len(tpl.DNSNames) > 0 || len(tpl.URIs) > 0 {
		tpl.ExtKeyUsage = append(tpl.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if len(tpl.EmailAddresses) > 0 {
		tpl.ExtKeyUsage = append(tpl.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}
	// IIS (the main target of PKCS #12 files), only shows the deprecated
	// Common Name in the UI. See issue #115.
	if m.pkcs12 {
		tpl.Subject.CommonName = hosts[0]
	}
	cert, err := x509.CreateCertificate(rand.Reader, tpl, m.caCert, pub, m.caKey)
	if err != nil {
		return errors.Wrapf(err, "generate certificate")
	}
	if !m.pkcs12 {
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
		m.certPEM = []byte{}
		m.certPEM = append(m.certPEM, certPEM...)
		privDER, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			return errors.Wrapf(err, "encode certifacate key")
		}
		privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})
		m.privPEM = []byte{}
		m.privPEM = append(m.privPEM, privPEM...)
	}
	return nil
}

func (m *CertGen) loadCert() error {
	//secondLvlWildcardRegexp := regexp.MustCompile(`(?i)^\*\.[0-9a-z_-]+$`)
	wildcardRegexp := regexp.MustCompile(`\+[0-9]+.pem`)
	secodeWildcardRegexp := regexp.MustCompile(`\+[0-9]+-key.pem`)
	files, err := listFilesInDir(m.CAROOT)
	if err != nil {
		return errors.Wrapf(err, "list cert files")
	}
	for _, file := range files {
		if wildcardRegexp.MatchString(file) {
			m.certPEM, _ = os.ReadFile(file)
		}
		if secodeWildcardRegexp.MatchString(file) {
			m.privPEM, _ = os.ReadFile(file)
		}
	}
	if len(m.certPEM) < 1 || len(m.privPEM) < 1 {
		return errors.New("load cert files failed")
	}
	return nil
}

func (m *CertGen) loadFromFS() error {
	d, err := certFS.ReadFile("cert/" + rootCaName)
	if err != nil {
		return errors.Wrapf(err, "read fs ca")
	}
	m.caPEM = []byte{}
	m.caPEM = append(m.caPEM, d...)
	d, err = certFS.ReadFile("cert/" + rootKeyName)
	if err != nil {
		return errors.Wrapf(err, "read fs ca key")
	}
	m.caKeyPEM = []byte{}
	m.caKeyPEM = append(m.caKeyPEM, d...)

	{
		certDERBlock, _ := pem.Decode(m.caPEM)
		if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
			return errors.New("read ca-cert unexpected content")
		}
		m.caCert, err = x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			return errors.Wrapf(err, "parse ca-cert")
		}

		keyDERBlock, _ := pem.Decode(m.caKeyPEM)
		if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
			return errors.New("read ca-key unexpected content")
		}
		m.caKey, err = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
		if err != nil {
			return errors.Wrapf(err, "parse ca-key")
		}
	}
	return nil
}

func (m *CertGen) loadCA() error {
	if !pathExists(filepath.Join(m.CAROOT, rootCaName)) {
		return errors.Errorf("%v file is missing", rootCaName)
	}
	certPEMBlock, err := os.ReadFile(filepath.Join(m.CAROOT, rootCaName))
	if err != nil {
		return errors.Wrapf(err, "read ca certificate")
	}
	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
		return errors.New("read ca-cert unexpected content")
	}
	m.caCert, err = x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return errors.Wrapf(err, "parse ca-cert")
	}
	if !pathExists(filepath.Join(m.CAROOT, rootKeyName)) {
		return errors.New("ca key file")
	}
	keyPEMBlock, err := os.ReadFile(filepath.Join(m.CAROOT, rootKeyName))
	if err != nil {
		return errors.Wrapf(err, "read ca key")
	}
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
		return errors.New("read ca-key unexpected content")
	}
	m.caKey, err = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		return errors.Wrapf(err, "parse ca-key")
	}
	return nil
}

/*
func (m *CertGen) newCA() error {
	priv, err := m.generateKey(true)
	if err != nil {
		return errors.Wrapf(err, "generate ca-key")
	}
	pub := priv.(crypto.Signer).Public()

	spkiASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return errors.Wrapf(err, "encode public key")
	}
	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return errors.Wrapf(err, "decode public key")
	}
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	sn, err := m.randomSerialNumber()
	if err != nil {
		return err
	}
	tpl := &x509.Certificate{
		Version:      3,
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization:       m.certOption.OrganizationCA,
			OrganizationalUnit: m.certOption.OrganizationalUnit,
			CommonName:         m.certOption.CommonNameCA,
			Country:            []string{"CN"},
			Province:           []string{"ZJ"},
			Locality:           []string{"HZ"},
		},
		SubjectKeyId:          skid[:],
		NotAfter:              time.Now().AddDate(0, 0, expiredDaysCA),
		NotBefore:             time.Now(),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, priv)
	if err != nil {
		return errors.Wrapf(err, "generate ca certificate")
	}
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return errors.Wrapf(err, "encode ca key")
	}
	err = os.WriteFile(filepath.Join(m.CAROOT, rootKeyName), pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}), 0400)
	if err != nil {
		return errors.Wrapf(err, "save ca-key")
	}
	err = os.WriteFile(filepath.Join(m.CAROOT, rootCaName), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert}), 0644)
	if err != nil {
		return errors.Wrapf(err, "save ca-certificate")
	}
	return nil
}*/

func (m *CertGen) generateKey(rootCA bool) (crypto.PrivateKey, error) {
	if m.ecdsa {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if rootCA {
		return rsa.GenerateKey(rand.Reader, 3072)
	}
	return rsa.GenerateKey(rand.Reader, 2048)
}

func (m *CertGen) randomSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "generate serial number")
	}
	return serialNumber, nil
}

func (m *CertGen) fileNames(hosts []string) (_certFile, _keyFile, _p12File string) {
	defaultName := strings.Replace(hosts[0], ":", "_", -1)
	defaultName = strings.Replace(defaultName, "*", "_wildcard", -1)
	if len(hosts) > 1 {
		defaultName += "+" + strconv.Itoa(len(hosts)-1)
	}
	if m.client {
		defaultName += "-client"
	}

	_certFile = filepath.Join(m.CAROOT, defaultName+".pem")
	if certFile != "" {
		_certFile = certFile
	}
	_keyFile = filepath.Join(m.CAROOT, defaultName+"-key.pem")
	if keyFile != "" {
		_keyFile = keyFile
	}
	_p12File = filepath.Join(m.CAROOT, defaultName+".p12")
	if p12File != "" {
		_p12File = p12File
	}
	return
}

func checkCertExpiry(pemData []byte) bool {
	block, _ := pem.Decode(pemData)
	if block == nil {
		fmt.Println("Failed to decode PEM data")
		return false
	}
	// parse X.509
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("Failed to parse certificate, %v\n", err)
		return false
	}
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return true
	} else {
		fmt.Println("The certificate is valid.")
		return false
	}
	fmt.Printf("Certificate is valid from %s to %s\n", cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339))
	return false
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getAllIps() ([]string, error) {
	var ips []string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.Wrapf(err, "get interface addresses")
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		if ip.To4() != nil {
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

func getCAROOT() string {
	if env := os.Getenv("CAROOT"); env != "" {
		return env
	}
	var dir string
	return filepath.Join(dir, ".certfiles")
}

func listFilesInDir(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
