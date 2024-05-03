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
