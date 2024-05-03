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
	logger1, err := c.NewLogger()
	assert.Nil(t, err)

	logger2, err := c.NewLogger()
	assert.Nil(t, err)

	fmt.Printf("logger1 [%v] logger2 [%v]", logger1, logger2)

	assert.Equal(t, logger1, logger2)

}
