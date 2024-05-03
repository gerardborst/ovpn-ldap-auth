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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFail(t *testing.T) {
	cn := &CNConfiguration{
		Fail: true,
	}
	ok, err := cn.Equal("User", "user")
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = cn.Equal("User", "otherUser")
	assert.False(t, ok)
	assert.Equal(t, "user [User], not equal to common name [otherUser] in client certificate", err.Error())
}

func TestWithOnlyWarn(t *testing.T) {
	cn := &CNConfiguration{
		Fail: false,
	}
	ok, err := cn.Equal("User", "otherUser")
	assert.True(t, ok)
	assert.Equal(t, "user [User], not equal to common name [otherUser] in client certificate", err.Error())
}
