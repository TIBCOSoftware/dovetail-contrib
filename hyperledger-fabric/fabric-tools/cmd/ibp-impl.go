// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func configNetwork(ibpConfig, outConfig, cryptoPath string) error {
	ibp, err := ioutil.ReadFile(ibpConfig)
	if err != nil {
		return err
	}
	var config ConnectionConfig
	if err = json.Unmarshal(ibp, &config); err != nil {
		return err
	}

	// configure channel-peers relationship
	setChannelConfig(&config)

	// write out tlsCert pem files, and update tlsCert file path in the config
	if err = setCryptoFiles(cryptoPath, &config); err != nil {
		return err
	}

	conn, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	writeFile(outConfig, conn)
	return nil
}

func setCryptoFiles(cryptoPath string, config *ConnectionConfig) error {
	// set client root crypto path
	config.Client.CryptoConfig = &CryptoConfig{Path: "${CRYPTO_PATH}"}

	// write orderer TLS certs
	for k, v := range config.Orderers {
		tlsPath := filepath.Join("orderers", k, "tlsca.pem")
		v.TLSCACerts.Path = filepath.Join("${CRYPTO_PATH}", tlsPath)
		if err := writeFile(filepath.Join(cryptoPath, tlsPath), []byte(v.TLSCACerts.PEM)); err != nil {
			return err
		}
	}

	// write peer and CA TLS certs for each org
	myorg := config.Client.Organization
	for k, v := range config.Organizations {
		// client user crypto path relative to the client cryptoPath
		userRoot := filepath.Join("organizations", v.MSPID, "users")
		if k == myorg {
			// create folder for client user crypto files
			userPath := Subst(filepath.Join(cryptoPath, userRoot))
			if err := os.MkdirAll(userPath, 0755); err != nil {
				return err
			}
		}
		v.CryptoPath = filepath.Join(userRoot, "{username}", "msp")
		for _, p := range v.Peers {
			if peer, ok := config.Peers[p]; ok {
				tlsPath := filepath.Join("peers", p, "tlsca.pem")
				peer.TLSCACerts.Path = filepath.Join("${CRYPTO_PATH}", tlsPath)
				if err := writeFile(filepath.Join(cryptoPath, tlsPath), []byte(peer.TLSCACerts.PEM)); err != nil {
					return err
				}
			}
		}
		for _, c := range v.CertificateAuthorities {
			if ca, ok := config.CertificateAuthorities[c]; ok {
				tlsPath := filepath.Join("certificateAuthorities", c, "tlsca.pem")
				ca.TLSCACerts.Path = filepath.Join("${CRYPTO_PATH}", tlsPath)
				if err := writeFile(filepath.Join(cryptoPath, tlsPath), []byte(ca.TLSCACerts.PEM)); err != nil {
					return err
				}
				if k == myorg {
					// print out TLS certificate path for myorg's CA
					fmt.Println(filepath.Join(cryptoPath, tlsPath))
				}
			}
		}
	}
	return nil
}

func setChannelConfig(config *ConnectionConfig) {
	myorg := config.Client.Organization
	mychannel := ChannelConfig{Peers: map[string]*PeerRoles{}}

	// collect peers owned by each org
	for k, v := range config.Organizations {
		peer := PeerRoles{
			EndorsingPeer:  true,
			ChaincodeQuery: true,
			LedgerQuery:    true,
			EventSource:    true,
		}
		if k != myorg {
			peer.ChaincodeQuery = false
			peer.EventSource = false
			peer.LedgerQuery = false
		}
		for _, p := range v.Peers {
			mychannel.Peers[p] = &peer
		}
	}
	config.Channels = map[string]*ChannelConfig{config.Name: &mychannel}
}

// ConnectionConfig defines attributes of fabric connection file exported from IBM Blockchain Platform
type ConnectionConfig struct {
	Name                   string                         `json:"name" yaml:"name"`
	Description            string                         `json:"description" yaml:"description"`
	Version                string                         `json:"version" yaml:"version"`
	Client                 *ClientConfig                  `json:"client" yaml:"client"`
	Channels               map[string]*ChannelConfig      `json:"-" yaml:"channels,omitempty"`
	Organizations          map[string]*OrganizationConfig `json:"organizations" yaml:"organizations"`
	Orderers               map[string]*OrdererConfig      `json:"orderers" yaml:"orderers"`
	Peers                  map[string]*PeerConfig         `json:"peers" yaml:"peers"`
	CertificateAuthorities map[string]*CAConfig           `json:"certificateAuthorities" yaml:"certificateAuthorities"`
}

// ClientConfig defines attributes of a fabric client
type ClientConfig struct {
	Organization string        `json:"organization" yaml:"organization"`
	CryptoConfig *CryptoConfig `json:"-" yaml:"cryptoconfig,omitempty"`
}

// ChannelConfig defines attributes of a channel
type ChannelConfig struct {
	Peers map[string]*PeerRoles `json:"-" yaml:"peers,omitempty"`
}

// PeerRoles defines roles of a peer for a channel
type PeerRoles struct {
	EndorsingPeer  bool `json:"-" yaml:"endorsingPeer"`
	ChaincodeQuery bool `json:"-" yaml:"chaincodeQuery"`
	LedgerQuery    bool `json:"-" yaml:"ledgerQuery"`
	EventSource    bool `json:"-" yaml:"eventSource"`
}

// CryptoConfig defines attributes of the client crypto root
type CryptoConfig struct {
	Path string `json:"-" yaml:"path,omitempty"`
}

// OrganizationConfig defines attributes of a fabric msp
type OrganizationConfig struct {
	MSPID                  string   `json:"mspid" yaml:"mspid"`
	CryptoPath             string   `json:"-" yaml:"cryptoPath,omitempty"`
	Peers                  []string `json:"peers" yaml:"peers"`
	CertificateAuthorities []string `json:"certificateAuthorities,omitempty" yaml:"certificateAuthorities,omitempty"`
}

// OrdererConfig defines attributes of a fabric orderer
type OrdererConfig struct {
	URL        string  `json:"url" yaml:"url"`
	TLSCACerts *CACert `json:"tlsCACerts" yaml:"tlsCACerts"`
}

// PeerConfig defines attriubtes of a fabric peer
type PeerConfig struct {
	URL         string   `json:"url" yaml:"url"`
	TLSCACerts  *CACert  `json:"tlsCACerts" yaml:"tlsCACerts"`
	GRPCOptions *GRPCOpt `json:"grpcOptions" yaml:"grpcOptions"`
}

// CAConfig defines attributes of a fabric CA server
type CAConfig struct {
	URL        string  `json:"url" yaml:"url"`
	CAName     string  `json:"caName" yaml:"caName"`
	TLSCACerts *CACert `json:"tlsCACerts" yaml:"tlsCACerts"`
}

// CACert defines attributes of a CA certificate
type CACert struct {
	PEM  string `json:"pem,omitempty" yaml:"-"`
	Path string `json:"-" yaml:"path,omitempty"`
}

// GRPCOpt defines attributes of GRPC options
type GRPCOpt struct {
	SSLTargetNameOverride string `json:"ssl-target-name-override" yaml:"ssl-target-name-override"`
}

func writeFile(path string, content []byte) error {
	p := Subst(path)
	d := filepath.Dir(p)
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(p, content, 0644)
}

// Subst replaces instances of '${VARNAME}' (eg ${GOPATH}) with the variable.
// Variables names that are not set by the SDK are replaced with the environment variable.
func Subst(path string) string {
	const (
		sepPrefix = "${"
		sepSuffix = "}"
	)

	splits := strings.Split(path, sepPrefix)

	var buffer bytes.Buffer

	// first split precedes the first sepPrefix so should always be written
	buffer.WriteString(splits[0]) // nolint: gas

	for _, s := range splits[1:] {
		subst, rest := substVar(s, sepPrefix, sepSuffix)
		buffer.WriteString(subst) // nolint: gas
		buffer.WriteString(rest)  // nolint: gas
	}

	return buffer.String()
}

// substVar searches for an instance of a variables name and replaces them with their value.
// The first return value is substituted portion of the string or noMatch if no replacement occurred.
// The second return value is the unconsumed portion of s.
func substVar(s string, noMatch string, sep string) (string, string) {
	endPos := strings.Index(s, sep)
	if endPos == -1 {
		return noMatch, s
	}

	v, ok := os.LookupEnv(s[:endPos])
	if !ok {
		return noMatch, s
	}
	return v, s[endPos+1:]
}
