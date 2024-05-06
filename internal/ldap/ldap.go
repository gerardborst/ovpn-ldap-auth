/*
OpenVPN ldap auth - OpenVPN Ldap authentication

Copyright (C) 2019 - 2021 Egbert Pot
Copyright (C) 2021 - 2024 Gerard Borst

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
func (lc *LDAPClient) Authenticate(username, password string) (authenticated bool) {
	authenticated = false
	// logger is already created with config in main
	logger = logging.GetLogger()

	//  https://github.com/go-ldap/ldap/issues/93
	if len(password) == 0 {
		logger.Error("zero length password not allowed", "username", username)
		return
	}

	conn, err := lc.connect()
	if err != nil {
		logger.Error("ldap connect error", "username", username, "error", err)
		return
	}
	defer conn.Close()

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		logger.Debug("Create connection with bind username / password", "username", username)
		err := conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			logger.Error("bind error", "error", err)
			return
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

	logger.Debug("search user in ldap", "searchRequest", searchRequest.Filter)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		logger.Error("ldap search error", "username", username, "error", err)
		return
	}

	if len(sr.Entries) < 1 {
		logger.Error("user does not exist, or is not a member of the OpenVPN group", "username", username)
		return
	}

	if len(sr.Entries) > 1 {
		logger.Error("too many entries returned", "username", username, "#entreis", len(sr.Entries))
		return
	}

	userDN := sr.Entries[0].DN
	logger.Debug("searched user", "userDN", userDN)

	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	// Bind as the user to verify their password
	err = conn.Bind(userDN, password)
	if err != nil {
		logger.Error("authentication error", "username", username, "error", err)
		return
	}
	authenticated = true
	logger.Info("Authentication successful", "username", username)
	return
}
