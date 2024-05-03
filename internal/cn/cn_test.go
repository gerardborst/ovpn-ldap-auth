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
