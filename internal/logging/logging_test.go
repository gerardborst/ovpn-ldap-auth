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

package logging_test

import (
	"fmt"
	"testing"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestToStdOut(t *testing.T) {
	c := logging.LogConfiguration{
		LogToFile: false,
	}
	logger1 := logging.NewLogger(&c)

	logger2 := logging.GetLogger()

	fmt.Printf("logger1 [%v] logger2 [%v]", logger1, logger2)

	assert.Equal(t, logger1, logger2)

}
