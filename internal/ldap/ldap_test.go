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

package ldap

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

var compose tc.ComposeStack
var ctx context.Context

func setup() {

	logger = logging.NewLogger(&logging.LogConfiguration{LogToFile: false})

	ctx = context.Background()
	var err error
	compose, err = tc.NewDockerCompose("../../tests/openldap/docker-compose.yaml")
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

func TestLdap(t *testing.T) {
	var tests = []struct {
		username string
		password string
		want     bool
	}{
		{"user01", "password1", true},  // success
		{"user01", "wrong", false},     // wrong password
		{"user01", "", false},          // empty password
		{"user02", "password2", false}, // not in group
	}
	for _, test := range tests {
		descr := fmt.Sprintf("Authenticate(%v, %v)",
			test.username, test.password)
		res := c1.Authenticate(test.username, test.password)
		if res != test.want {
			t.Errorf("%s = %v, want %v", descr, res, test.want)
		}
	}
}
