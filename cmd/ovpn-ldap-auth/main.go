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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gerardborst/ovpn-ldap-auth/internal/cn"
	"github.com/gerardborst/ovpn-ldap-auth/internal/ldap"
	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
	"github.com/gerardborst/ovpn-ldap-auth/internal/report"
	"github.com/spf13/viper"
)

var (
	CommitHash  string
	VersionTag  = "DEVELOPMENT"
	BuildTime   string
	showVersion = flag.Bool("v", false, "show version information")
)

type Configuration struct {
	LdapClient ldap.LDAPClient
	Log        logging.LogConfiguration
	CN         cn.CNConfiguration
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("ovpn-ldap-auth version: [%s] commit: [%s] build time: [%s]\n", VersionTag, CommitHash, BuildTime)
		os.Exit(0)
	}
	viper.SetConfigName("ovpn-auth-config")      // name of config file (without extension)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                     // optionally look for config in the working directory
	viper.AddConfigPath("../../tests/openldap/") // optionally look in the tests directory
	viper.AddConfigPath("./tests/openldap/")     // optionally look in the tests directory
	viper.AddConfigPath("/etc/openvpn/auth/")    // path to look for the config file in

	viper.SetDefault("ldapClient.UseSSL", true)
	viper.SetDefault("ldapClient.SkipTLS", false)

	viper.SetDefault("log.Level", "info")
	viper.SetDefault("log.File", "/var/log/openvpn/auth/ldap-auth.log")
	viper.SetDefault("log.LogToFile", false)

	viper.SetDefault("CN.Check", true)
	viper.SetDefault("CN.Fail", true)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("fatal error reading config file: %v", err)
	}

	var c Configuration

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logger := logging.NewLogger(&c.Log)
	if err != nil {
		log.Fatalf("unable initialize logger, %v", err)
	}

	viper.BindEnv("username", "username")
	viper.BindEnv("password", "password")
	viper.BindEnv("common_name", "common_name")
	viper.BindEnv("auth_control_file", "auth_control_file")

	username := viper.GetString("username")
	password := viper.GetString("password")
	commonName := viper.GetString("common_name")
	authControlFile, err := os.OpenFile(viper.GetString("auth_control_file"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Error("Open auth_control_file", "error", err)
		os.Exit(1)
	}

	reporter := report.NewReporter(authControlFile)

	logger.Debug("", "configuration", c)

	// Check common name in client certificate
	abort := c.CN.CheckCN(username, commonName)
	if abort {
		reporter.Report(false)
		logger.Error("Authentication not successful", "username", username)
		return
	}
	// Ldap Authenticate
	authenticated := c.LdapClient.Authenticate(username, password)
	reporter.Report(authenticated)
	if !authenticated {
		logger.Error("Authentication not successful", "username", username)
	}
}
