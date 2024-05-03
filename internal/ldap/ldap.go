// Package ldap provides a simple ldap client to authenticate,
// retrieve basic information and groups for a user.
package ldap

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
	"gopkg.in/ldap.v2"
)

type LDAPClient struct {
	Attributes         []string
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string // e.g. "(memberUid=%s)"
	Host               string
	ServerName         string
	VpnGroupFilter     string // e.g. "(uid=%s)"
	Port               int
	InsecureSkipVerify bool
	UseSSL             bool
	UseStartTls        bool
	ClientCertificates []tls.Certificate // Adding client certificates
}

var logger *slog.Logger

var conn *ldap.Conn

// Connect connects to the ldap backend.
func (lc *LDAPClient) connect() (*ldap.Conn, error) {
	var err error
	address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
	if !lc.UseSSL {
		logger.Debug("Connecting WITHOUT TLS")
		conn, err = ldap.Dial("tcp", address)
		if err != nil {
			return nil, err
		}
		if lc.UseStartTls {
			// Reconnect with TLS
			err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		logger.Debug("Connecting with TLS")
		config := &tls.Config{
			InsecureSkipVerify: lc.InsecureSkipVerify,
			ServerName:         lc.ServerName,
		}
		if lc.ClientCertificates != nil && len(lc.ClientCertificates) > 0 {
			config.Certificates = lc.ClientCertificates
		}
		conn, err = ldap.DialTLS("tcp", address, config)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
	var err error
	// logger is already created with config in main
	lgc := logging.LogConfiguration{}
	logger, err = lgc.NewLogger()
	if err != nil {
		log.Fatalf("unable to initialize logger, %v", err)
	}
	//  https://github.com/go-ldap/ldap/issues/93
	if len(password) == 0 {
		return false, nil, fmt.Errorf("zero length password not allowed, user [%v]", username)
	}

	conn, err := lc.connect()
	if err != nil {
		return false, nil, err
	}
	defer conn.Close()

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		logger.Debug("Create connection with bind username / password")
		err := conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return false, nil, err
		}
		logger.Debug("Connection with bind account successful")
	}

	attributes := append(lc.Attributes, "dn")
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 1, false,
		fmt.Sprintf(lc.VpnGroupFilter, username),
		attributes,
		nil,
	)

	logger.Debug("", "searchRequest", searchRequest.Filter)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) < 1 {
		return false, nil, fmt.Errorf("user [%s] does not exist, or is not a member of the OpenVPN group", username)
	}

	if len(sr.Entries) > 1 {
		return false, nil, fmt.Errorf("too many entries returned: %v", len(sr.Entries))
	}

	userDN := sr.Entries[0].DN
	logger.Debug("", "userDN", userDN)

	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	// Bind as the user to verify their password
	err = conn.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	return true, user, nil
}
