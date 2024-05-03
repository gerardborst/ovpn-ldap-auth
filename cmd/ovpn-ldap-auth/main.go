package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gerardborst/ovpn-ldap-auth/internal/cn"
	"github.com/gerardborst/ovpn-ldap-auth/internal/ldap"
	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
	"github.com/spf13/viper"
)

var (
	CommitHash string
	VersionTag string
)

type Configuration struct {
	LdapClient ldap.LDAPClient
	Log        logging.LogConfiguration
	CN         cn.CNConfiguration
}

var c Configuration

var username, authControlFile string

var logger *slog.Logger

func main() {
	viper.SetConfigName("ovpn-auth-config")              // name of config file (without extension)
	viper.SetConfigType("yaml")                          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                             // optionally look for config in the working directory
	viper.AddConfigPath("../../tests/openldap/openldap") // optionally look in the tests directory
	viper.AddConfigPath("/etc/openvpn/auth/")            // path to look for the config file in

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

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logger, err = c.Log.NewLogger()
	if err != nil {
		log.Fatalf("unable initialize logger, %v", err)
	}

	viper.BindEnv("username", "username")
	viper.BindEnv("password", "password")
	viper.BindEnv("common_name", "common_name")
	viper.BindEnv("auth_control_file", "auth_control_file")

	username = viper.GetString("username")
	password := viper.GetString("password")
	commonName := viper.GetString("common_name")
	authControlFile = viper.GetString("auth_control_file")

	logger.Info("ldap authentication", "version", VersionTag, "commit", CommitHash, "username", username)

	logger.Debug("", "configuration", c)

	// Check common name in clietn certificate
	if c.CN.Check {
		ok, err := c.CN.Equal(username, commonName)
		if err != nil {
			logger.Error(err.Error())
		}
		if !ok {
			reportSuccess(ok)
			return
		}
	}
	// Ldap Authenticate
	ok, user, err := c.LdapClient.Authenticate(username, password)
	if err != nil {
		logger.Error("Authentication errored", "username", username, "error", err)
		reportSuccess(false)
	} else {
		if !ok {
			logger.Error("Authentication failed", "username", username, "error", err)
			reportSuccess(false)
		} else {
			logger.Info("Authentication successful", "user", user)
			reportSuccess(true)
		}
	}
}

func reportSuccess(authSuccess bool) {
	if authSuccess {
		err := os.WriteFile(authControlFile, []byte("1"), 0644)
		if err != nil {
			logger.Error("WriteFile errored for user %s, error: %v", username, err)
		}
	}
	err := os.WriteFile(authControlFile, []byte("0"), 0644)
	if err != nil {
		logger.Error("WriteFile errored for user %s, error: %v", username, err)
	}

}
