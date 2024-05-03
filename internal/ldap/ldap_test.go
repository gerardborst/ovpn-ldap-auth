package ldap

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

var compose tc.ComposeStack
var ctx context.Context

func setup() {
	ctx = context.Background()
	var err error
	compose, err = tc.NewDockerCompose("../../tests/openldap/openldap/docker-compose.yaml")
	if err != nil {
		panic(err)
	}

	err = compose.Up(ctx, tc.Wait(true))
	if err != nil {
		panic(err)
	}

}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()
	os.Exit(code)
}

func tearDown() {
	compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal)
}

var c1 = &LDAPClient{
	Base:           "dc=example,dc=org",
	GroupFilter:    "(memberOf=%s)",
	Host:           "localhost",
	Port:           1389,
	UseSSL:         false,
	UseStartTls:    false,
	VpnGroupFilter: "(&(uid=%s)(memberOf=cn=openvpn,ou=users,dc=example,dc=org))",
	ServerName:     "localhost",
	BindDN:         "cn=admin,dc=example,dc=org",
	BindPassword:   "123456",
}

func TestWithOutTLS(t *testing.T) {
	res, users, err := c1.Authenticate("user01", "password1")
	assert.Nil(t, err)
	assert.True(t, res, "Authentication should be succesfull")
	assert.Equal(t, map[string]string{}, users)

	res, users, err = c1.Authenticate("user01", "zomaar")
	assert.Equal(t, "LDAP Result Code 49 \"Invalid Credentials\": ", err.Error())
	assert.False(t, res)
	assert.Equal(t, map[string]string{}, users)

	res, users, err = c1.Authenticate("user01", "")
	assert.Equal(t, "zero length password not allowed, user [user01]", err.Error())
	assert.False(t, res)
	assert.Equal(t, map[string]string(nil), users)

	res, users, err = c1.Authenticate("user02", "password2")
	assert.Equal(t, "user [user02] does not exist, or is not a member of the OpenVPN group", err.Error())
	assert.False(t, res, "Authentication should be not succesfull")
	assert.Nil(t, users, "users should be nil")
}
