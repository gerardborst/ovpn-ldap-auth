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

package cn

import (
	"fmt"
	"testing"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
)

func TestReport(t *testing.T) {

	logger = logging.NewLogger(&logging.LogConfiguration{LogToFile: false})

	var cn *CNConfiguration
	var tests = []struct {
		check bool
		fail  bool
		user  string
		cn    string
		abort bool
	}{
		{false, false, "user", "user", false},
		{false, true, "user", "user", false},
		{true, false, "user", "user", false},
		{true, false, "user", "otherUser", false},
		{true, true, "user", "user", false},
		{true, true, "user", "User", false},
		{true, true, "user", "otherUser", true},
	}
	for _, test := range tests {
		cn = &CNConfiguration{
			Check: test.check,
			Fail:  test.fail,
		}
		descr := fmt.Sprintf("CheckCN(%v, %v) CNConfiguration{Check: %v, Fail:  %v}",
			test.user, test.cn, cn.Check, cn.Fail)
		abort := cn.CheckCN(test.user, test.cn)

		if test.abort != abort {
			t.Errorf("%s abort = %v, want %v", descr, abort, test.abort)
		}
	}
}
