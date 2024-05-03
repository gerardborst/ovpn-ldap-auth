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
	"strings"
)

type CNConfiguration struct {
	Check bool
	Fail  bool
}

func (cn *CNConfiguration) Equal(username, commonName string) (bool, error) {
	if !strings.EqualFold(username, commonName) {
		err := fmt.Errorf("user [%s], not equal to common name [%s] in client certificate", username, commonName)
		if cn.Fail {
			return false, err
		} else {
			return true, err
		}
	}
	return true, nil
}
